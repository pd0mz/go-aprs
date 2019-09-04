package aprs

import (
	"math"
	"testing"
	"time"
)

const earthRadius = float64(6378100)

func testTime(day, hour, min, sec int) *time.Time {
	t := time.Date(0, 0, day, hour, min, sec, 0, time.UTC)
	return &t
}

// haversin(θ) function
func testHaversin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func testDistance(a *Position, b *Position) float64 {
	var (
		lat1 = a.Latitude * math.Pi / 180
		lng1 = a.Longitude * math.Pi / 180
		lat2 = b.Latitude * math.Pi / 180
		lng2 = b.Longitude * math.Pi / 180
		h    = testHaversin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*testHaversin(lng2-lng1)
	)
	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}

func TestPacket(t *testing.T) {
	var tests = []struct {
		Raw      string
		Src      *Address
		Dst      *Address
		PathLen  int
		Type     DataType
		Position *Position
		Velocity *Velocity
		PHG      *PowerHeightGain
		DFS      *OmniDFStrength
		Altitude float64
		Range    float64
		Time     *time.Time
	}{
		{
			Raw:      "N0CALL>APRS,qAC:!4903.50N/07201.75W-Test 001234",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('!'),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
		},
		{
			Raw:      "N0CALL>APRS,qAC:!4903.50N/07201.75W-Test /A=001234",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('!'),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
		},
		{
			Raw:      "N0CALL>APRS,qAC:!49  .  N/072  .  W-",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('!'),
			Position: &Position{Latitude: 49.0, Longitude: -72.000000},
		},
		/*
			{
				Raw:      "N0CALL>APRS,qAC:TheNet X-1J4  (BFLD)!4903.50N/07201.75Wn",
				Src:      MustParseAddress("N0CALL"),
				Dst:      MustParseAddress("APRS"),
				PathLen:  1,
				Type:     DataType('!'),
				Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
			},
		*/
		{
			Raw:      "N0CALL>APRS,qAC:@092345/4903.50N/07201.75W>Test1234",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('@'),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
			Time:     testTime(9, 23, 45, 0),
		},
		{
			Raw:      "N0CALL>APRS,qAC:=4903.50N/07201.75W#PHG5132",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('='),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
			PHG:      &PowerHeightGain{'5', '1', '3', '2'},
		},
		{
			Raw:      "N0CALL>APRS,qAC:=4903.50N/07201.75W 225/000g000t050r000p001...h00b10138dU2k",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('='),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
		},
		{
			Raw:      "N0CALL>APRS,qAC:@092345/4903.50N/07201.75W>088/036",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('@'),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
			Time:     testTime(9, 23, 45, 0),
		},
		{
			Raw:      "N0CALL>APRS,qAC:@234517h4903.50N/07201.75W>PHG5132",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('@'),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
			Time:     testTime(0, 23, 45, 17),
			PHG:      &PowerHeightGain{'5', '1', '3', '2'},
		},
		{
			Raw:      "N0CALL>APRS,qAC:@092345z4903.50N/07201.75W>RNG0050",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('@'),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
			Time:     testTime(9, 23, 45, 0),
			Range:    50,
		},
		{
			Raw:      "N0CALL>APRS,qAC:/234517h4903.50N/07201.75W>DFS2360",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('/'),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
			Time:     testTime(0, 23, 45, 17),
			DFS:      &OmniDFStrength{'2', '3', '6', '0'},
		},
		{
			Raw:      "N0CALL>APRS,qAC:@092345z4903.50N/07201.75W 090/000g000t066r000p000...dUII",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('@'),
			Position: &Position{Latitude: 49.058333, Longitude: -72.029167},
			Time:     testTime(9, 23, 45, 00),
		},
		{
			Raw:      "N0CALL>APRS,qAC:[IO91SX] 35 miles NNW of London",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('['),
			Position: &Position{Latitude: 51.958333, Longitude: -0.500000},
		},
		{
			Raw:      "N0CALL>APRS,qAC:[IO91]",
			Src:      MustParseAddress("N0CALL"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  1,
			Type:     DataType('['),
			Position: &Position{Latitude: 51.0, Longitude: -2.0},
		},
		{
			Raw:      "WX4GSO-9>APN382,qAR,WD4LSS:!3545.18NL07957.08W#PHG5680/R,W,85NC,NCn Mount Shepherd Piedmont Triad NC",
			Src:      MustParseAddress("WX4GSO-9"),
			Dst:      MustParseAddress("APN382"),
			PathLen:  2,
			Type:     DataType('!'),
			Position: &Position{Latitude: 35.753000, Longitude: -79.951333},
		},
		{
			Raw:      "PA4TW>APRS,qAS,PA4TW-2:=5216.28N/00510.05Er Remco, DMR:2041014, Soms QRV op PI2NOS",
			Src:      MustParseAddress("PA4TW"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  2,
			Type:     DataType('='),
			Position: &Position{Latitude: 52.271333, Longitude: 5.167500},
		},
		{
			Raw:      "PA4TW-10>APRS,TCPIP*,qAC,FOURTH:=5220.18N/00453.25EIhttp://aprs.pa4tw.nl:14501/",
			Src:      MustParseAddress("PA4TW-10"),
			Dst:      MustParseAddress("APRS"),
			PathLen:  3,
			Type:     DataType('='),
			Position: &Position{Latitude: 52.336333, Longitude: 4.887500},
		},
		{
			Raw:      "N0CALL-1>T3PY1Y,KQ1L-8*,WIDE1,WIDE2-1,qAR:=/5L!!<*e7> sTComment",
			Src:      MustParseAddress("N0CALL-1"),
			Dst:      MustParseAddress("T3PY1Y"),
			PathLen:  4,
			Type:     DataType('='),
			Position: &Position{Latitude: 49.5, Longitude: -72.75},
		},
		{
			Raw:      "N0CALL-1>T3PY1Y,KQ1L-8*,WIDE1,WIDE2-1,qAR:=/5L!!<*e7>7P[",
			Src:      MustParseAddress("N0CALL-1"),
			Dst:      MustParseAddress("T3PY1Y"),
			PathLen:  4,
			Type:     DataType('='),
			Position: &Position{Latitude: 49.5, Longitude: -72.75},
			Velocity: &Velocity{88.0, 36.2},
		},
		{
			Raw:      "N0CALL-1>T3PY1Y,KQ1L-8*,WIDE1,WIDE2-1,qAR:=/5L!!<*e7>{?!",
			Src:      MustParseAddress("N0CALL-1"),
			Dst:      MustParseAddress("T3PY1Y"),
			PathLen:  4,
			Type:     DataType('='),
			Position: &Position{Latitude: 49.5, Longitude: -72.75},
			Range:    20.13,
		},
		{
			Raw:      "N0CALL-1>T3PY1Y,KQ1L-8*,WIDE1,WIDE2-1,qAR:=/5L!!<*e7OS]S",
			Src:      MustParseAddress("N0CALL-1"),
			Dst:      MustParseAddress("T3PY1Y"),
			PathLen:  4,
			Type:     DataType('='),
			Position: &Position{Latitude: 49.5, Longitude: -72.75},
			Altitude: 10004,
		},
		{
			Raw:      "N0CALL-1>T3PY1Y,KQ1L-8*,WIDE1,WIDE2-1,qAR:@092345z/5L!!<*e7>{?!",
			Src:      MustParseAddress("N0CALL-1"),
			Dst:      MustParseAddress("T3PY1Y"),
			PathLen:  4,
			Type:     DataType('@'),
			Time:     testTime(9, 23, 45, 0),
			Position: &Position{Latitude: 49.5, Longitude: -72.75},
			Range:    20.13,
		},
	}

	for _, test := range tests {
		p, err := ParsePacket(test.Raw)
		if err != nil {
			panic(err)
		}
		if !p.Dst.EqualTo(test.Dst) {
			t.Fatalf("expected dst %s, got %s", test.Dst, p.Dst)
		}
		if !p.Src.EqualTo(test.Src) {
			t.Fatalf("expected src %s, got %s", test.Src, p.Src)
		}
		if len(p.Path) != test.PathLen {
			t.Fatalf("expected path length %d, got %d: %v", test.PathLen, len(p.Path), p.Path)
		}
		if p.Payload.Type() != test.Type {
			t.Fatalf("expected packet type %s [%c], got %s [%c]",
				test.Type, test.Type,
				p.Payload.Type(), p.Payload.Type())
		}
		if test.Altitude != 0 {
			if p.Altitude == 0 {
				t.Fatalf("expected altitude %f, got none", test.Altitude)
			}
			if math.Abs(test.Altitude-p.Altitude) > 1.0 {
				t.Fatalf("expected altitude %f, got %f", test.Altitude, p.Altitude)
			}
		}
		if test.Velocity != nil {
			if p.Velocity.Course == 0 {
				t.Fatalf("expected velocity %v, got none", test.Velocity)
			}
			if math.Abs(test.Velocity.Course-p.Velocity.Course) > 1.0 {
				t.Fatalf("expected course %f, got %f", test.Velocity.Course, p.Velocity.Course)
			}
			if math.Abs(test.Velocity.Speed-p.Velocity.Speed) > 1.0 {
				t.Fatalf("expected speed %f, got %f", test.Velocity.Speed, p.Velocity.Speed)
			}
		}
		if test.Range != 0 {
			if p.Range == 0 {
				t.Fatalf("expected range %f, got none", test.Range)
			}
			if math.Abs(test.Range-p.Range) > 0.1 {
				t.Fatalf("expected range %f, got %f", test.Range, p.Range)
			}
		}
		if test.Position != nil {
			if p.Position == nil {
				t.Fatalf("expected position %s, got none", test.Position)
			}
			if d := testDistance(test.Position, p.Position); d > 1.0 {
				t.Fatalf("expected position %s, got %s with distance %f meter", test.Position, p.Position, d)
			}
		}
		if test.Time != nil {
			if p.Time == nil {
				t.Fatalf("expected time %s", test.Time)
			}
			if test.Time.Sub(*p.Time) > time.Minute {
				t.Fatalf("expected time %s, got %s", test.Time, p.Time)
			}
		}
	}
}
