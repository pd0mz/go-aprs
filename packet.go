package aprs

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Payload string

func (p Payload) Type() DataType {
	var t DataType

	if len(p) > 0 {
		t = DataType(p[0])
	}

	return t
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (p Payload) Len() int { return len(p) }

type Velocity struct {
	Course float64 // Degrees
	Speed  float64 // Knots
}

type Wind struct {
	Direction float64 // Degrees
	Speed     float64 // Knots
}

type PowerHeightGain struct {
	PowerCode       byte
	HeightCode      byte
	GainCode        byte
	DirectivityCode byte
}

func (p PowerHeightGain) Power() int {
	w := int(p.PowerCode - '0')
	if w <= 0 {
		return 0
	}
	return w * w
}

func (p PowerHeightGain) Height() float64 {
	h := float64(p.HeightCode - '0')
	if h <= 0 {
		return 10
	}
	return math.Pow(2, h) * 10
}

func (p PowerHeightGain) Gain() int {
	d := int(p.GainCode - '0')
	if d <= 0 {
		return 0
	}
	return d
}

func (p PowerHeightGain) Directivity() float64 {
	d := int(p.DirectivityCode - '0')
	if d <= 0 {
		return 0
	}
	return float64(d%8) * 45.0
}

type OmniDFStrength struct {
	StrengthCode    byte
	HeightCode      byte
	GainCode        byte
	DirectivityCode byte
}

func (o OmniDFStrength) Strength() int {
	w := int(o.StrengthCode - '0')
	if w <= 0 {
		return 0
	}
	return w * w
}

func (o OmniDFStrength) Height() float64 {
	h := float64(o.HeightCode - '0')
	if h <= 0 {
		return 10
	}
	return math.Pow(2, h) * 10
}

func (o OmniDFStrength) Gain() int {
	d := int(o.GainCode - '0')
	if d <= 0 {
		return 0
	}
	return d
}

func (o OmniDFStrength) Directivity() float64 {
	d := int(o.DirectivityCode - '0')
	if d <= 0 {
		return 0
	}
	return float64(d%8) * 45.0
}

type Packet struct {
	Raw      string
	Src      *Address
	Dst      *Address
	Path     Path
	Payload  Payload
	Position *Position
	Time     *time.Time
	Altitude float64 // Feet
	Velocity Velocity
	Wind     Wind
	PHG      PowerHeightGain
	DFS      OmniDFStrength
	Range    float64 // Miles
	Comment  string
	data     string // Unparsed data
}

func ParsePacket(raw string) (p Packet, err error) {
	p = Packet{Raw: raw}

	var i int
	if i = strings.Index(raw, ":"); i < 0 {
		return p, errors.New("aprs: invalid packet")
	}
	p.Payload = Payload(raw[i+1:])

	// Parse src, dst and path
	var a = raw[:i]
	if i = strings.Index(a, ">"); i < 0 {
		return p, errors.New("aprs: invalid packet")
	}
	if p.Src, err = ParseAddress(a[:i]); err != nil {
		return
	}
	var r = strings.Split(a[i+1:], ",")
	if p.Dst, err = ParseAddress(r[0]); err != nil {
		return
	}
	if p.Path, err = ParsePath(strings.Join(r[1:], ",")); err != nil {
		return
	}

	// Post processing of payload
	err = p.parse()

	return
}

func (p *Packet) parse() error {
	s := string(p.Payload)
	//log.Printf("parse %q [%c]\n", s, p.Payload.Type())

	switch p.Payload.Type() {
	case '!': // Lat/Long Position Report Format — without Timestamp
		var o = strings.IndexByte(s, '!')
		pos, txt, err := ParsePosition(s[o+1:], !isDigit(s[o+1]))
		//log.Printf("parse result: %s, %q, %v\n", pos, txt, err)
		if err != nil {
			return err
		}
		p.Position = &pos
		p.data = txt
	case '=':
		compressed := IsValidCompressedSymTable(s[1])
		pos, txt, err := ParsePosition(s[1:], compressed)
		if err != nil {
			return err
		}
		p.Position = &pos
		p.data = txt
	case '/', '@': // Lat/Long Position Report Format — with Timestamp
		if len(s) < 8 {
			return errors.New("aprs: invalid position")
		}

		if s[7] == 'h' || s[7] == 'z' || s[7] == '/' {
			if ts, err := ParseTime(s[1:]); err == nil {
				p.Time = &ts
			}
			compressed := IsValidCompressedSymTable(s[8])
			pos, txt, err := ParsePosition(s[8:], compressed)
			if err != nil {
				return err
			}
			p.Position = &pos
			p.data = txt
		} else if s[7] >= '0' && s[7] <= '9' {
			ts, err := ParseTime(s[1:])
			if err != nil {
				return err
			}
			p.Time = &ts
			compressed := IsValidCompressedSymTable(s[10])
			pos, txt, err := ParsePosition(s[10:], compressed)
			if err != nil {
				return err
			}
			p.Position = &pos
			p.data = txt
		}
	case ';':
		pos, txt, err := ParsePosition(s[18:], !isDigit(s[18]))
		if err != nil {
			return err
		}
		p.Position = &pos
		p.data = txt
	case '[':
		pos, txt, err := ParsePositionGrid(s[1:])
		if err != nil {
			return err
		}
		p.Position = &pos
		p.data = txt
	case '`', '\'':
		pos, err := ParseMicE(s, p.Dst.Call)
		if err != nil {
			return err
		}
		p.Position = &pos
		p.parseMicEData()

		return nil // there is no additional data to parse
	default:
		pos, txt, err := ParsePositionBoth(s)
		if err != nil {
			return err
		}
		p.Position = &pos
		p.data = txt
	}

	if p.Position != nil {
		if p.Position.Compressed {
			return p.parseCompressedData()
		}
		return p.parseData()
	}

	return nil
}

func (p *Packet) parseMicEData() error {
	// APRS PROTOCOL REFERENCE 1.0.1 Chapter 10, page 42 in PDF
	s := string(p.Payload)

	// Mic-E Message Type
	var mt []string
	var t string
	for i := 0; i < 3; i++ {
		mc := miceCodes[rune(p.Dst.Call[i])][1]
		if strings.HasSuffix(mc, "(Custom)") {
			t = messageTypeCustom
		} else if strings.HasSuffix(mc, "(Std)") {
			t = messageTypeStd
		}
		mt = append(mt, string(mc[0]))
	}
	switch t {
	case messageTypeStd:
		mt = append(mt, " (Std)")
	case messageTypeCustom:
		mt = append(mt, " (Custom)")
	}
	p.Comment = miceMsgTypes[strings.Join(mt, "")]

	// Speed and Course.
	speed := float64(int(s[4])-28) * 10
	dc := float64(int(s[5])-28) / 10
	unit := float64(int(dc))
	speed += unit
	course := dc - unit
	course += float64(int(s[6]) - 28)
	if speed >= 800 {
		speed -= 800
	}
	speed = 1.852 * speed // convert speed from knots to km/h
	if course >= 400 {
		course -= 400
	}

	p.Comment += fmt.Sprintf(" (%.fkm/h, %.f°)", speed, course)

	// Check whether there's additional Telemetry or Status Text data.
	if len(s) == 9 {
		return nil
	}

	if s[9] == ',' || s[9] == '\x1d' {
		// TODO: Parse telemetry data.
		return nil
	}

	// Parse MicE Status Text data.
	// TODO: Parse additional data in the Status Text data:
	// - Actual (custom) text message
	// - Maidenhead locator
	// - Altitude

	return nil
}

func (p *Packet) parseCompressedData() error {
	// Parse csT bytes
	if len(p.data) >= 3 {
		// Compression Type (T) Byte Format
		// Bit: 7      | 6      | 5       | 4     3     | 2    1    0      |
		//	-------+--------+---------+-------------+------------------+
		//      Unused | Unused | GPS Fix | NMEA Source | Origin           |
		//	-------+--------+---------+-------------+------------------+
		// Val: 0      | 0      | 0 = old | 00 = other  | 000 = Compressed |
		//	       |        | 1 = cur | 01 = GLL    | 001 = TNC BTex   |
		//	       |        |         | 10 = CGA    | 010 = Software   |
		//	       |        |         | 11 = RMC    | 011 = [tbd]      |
		//	       |        |         |             | 100 = KPC3       |
		//	       |        |         |             | 101 = Pico       |
		//	       |        |         |             | 110 = Other      |
		//	       |        |         |             | 111 = Digipeater |
		cb := p.data[0] - 33
		sb := p.data[1] - 33
		Tb := p.data[2] - 33
		if p.data[0] != ' ' && ((Tb>>3)&3) == 2 {
			// CGA sentence, NMEA Source = 0b10
			d, err := base91Decode(p.data[0:2])
			if err != nil {
				return err
			}
			p.Altitude = math.Pow(1.002, float64(d))
			p.Comment = p.data[3:]
		} else if cb >= 0 && cb <= 89 { // !..z
			// Course/Speed
			p.Velocity.Course = float64(cb) * 4.0
			p.Velocity.Speed = math.Pow(1.08, float64(sb)) - 1.0
		} else if cb == 90 { // {
			// Pre-Calculated Radio Range
			p.Range = 2 * math.Pow(1.08, float64(sb))
		}
	}
	return nil
}

func (p *Packet) parseData() error {
	//log.Printf("data %q\n", p.data)
	switch {
	case len(p.data) >= 1 && p.data[0] == ' ':
		p.Comment = p.data[1:]

	case len(p.data) >= 7 && strings.HasPrefix(p.data, "PHG"):
		p.PHG.PowerCode = p.data[3]
		p.PHG.HeightCode = p.data[4]
		p.PHG.GainCode = p.data[5]
		p.PHG.DirectivityCode = p.data[6]
		p.Range = math.Sqrt(2 * p.PHG.Height() * math.Sqrt((float64(p.PHG.Power())/10)*(float64(p.PHG.Gain())/2)))
		p.Comment = p.data[7:]

	case len(p.data) >= 7 && strings.HasPrefix(p.data, "RNG"):
		var err error
		p.Range, err = strconv.ParseFloat(p.data[3:7], 64)
		if err != nil {
			return err
		}
		p.Comment = p.data[7:]

	case len(p.data) >= 7 && strings.HasPrefix(p.data, "DFS"):
		p.DFS.StrengthCode = p.data[3]
		p.DFS.HeightCode = p.data[4]
		p.DFS.GainCode = p.data[5]
		p.DFS.DirectivityCode = p.data[6]
		p.Comment = p.data[7:]
	}
	return nil
}

func (p Payload) Time() (time.Time, error) {
	switch p.Type() {
	case '/', '@':
		return ParseTime(string(p)[1:])
	default:
		return time.Time{}, nil
	}
}
