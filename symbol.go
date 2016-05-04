package aprs

import "fmt"

type Symbol [2]byte

func (s Symbol) IsPrimaryTable() bool { return s[0] != '\\' }

func (s Symbol) String() string {
	var m map[byte]string
	if s.IsPrimaryTable() {
		m = primarySymbol
	} else {
		m = alternateSymbol
	}
	if n, ok := m[s[1]]; ok {
		return n
	}
	return fmt.Sprintf("unknown symbol %c", s[1])
}

func IsValidCompressedSymTable(c byte) bool {
	return c == '/' ||
		c == '\\' ||
		(c >= 0x41 && c <= 0x5a) ||
		(c >= 0x61 && c <= 0x6a)
}

func IsValidUncompressedSymTable(c byte) bool {
	return c == '/' ||
		c == '\\' ||
		(c >= 0x41 && c <= 0x5a) ||
		(c >= 0x30 && c <= 0x39)
}

// http://www.aprs.org/symbols/symbolsX.txt
var (
	primarySymbol = map[byte]string{
		'!':  "Police Station",
		'"':  "",
		'#':  "Digi",
		'$':  "Phone",
		'%':  "DX Cluster",
		'&':  "HF Gateway",
		'\'': "Plane",
		'(':  "Mobile Satellite Station",
		')':  "Wheelchair",
		'*':  "Snowmobile",
		'+':  "Red Cross",
		',':  "Boy Scout",
		'-':  "Home",
		'.':  "X",
		'/':  "Red Dot",
		'0':  "Circle (0)",
		'1':  "Circle (1)",
		'2':  "Circle (2)",
		'3':  "Circle (3)",
		'4':  "Circle (4)",
		'5':  "Circle (5)",
		'6':  "Circle (6)",
		'7':  "Circle (7)",
		'8':  "Circle (8)",
		'9':  "Circle (9)",
		':':  "Fire",
		';':  "Campground",
		'<':  "Motorcycle",
		'=':  "Rail engine",
		'>':  "Car",
		'?':  "File Server",
		'@':  "HC Future",
		'A':  "Aid Station",
		'B':  "BBS",
		'C':  "Canoe",
		'D':  "",
		'E':  "Eyeball",
		'F':  "Tractor",
		'G':  "Grid Square",
		'H':  "Hotel",
		'I':  "TCP/IP",
		'J':  "",
		'K':  "School",
		'L':  "User Log-ON",
		'M':  "MacAPRS",
		'N':  "NTS Station",
		'O':  "Balloon",
		'P':  "Police",
	}
	alternateSymbol = map[byte]string{
		'!':  "Emergency",
		'"':  "",
		'#':  "No. Digi",
		'$':  "Bank",
		'%':  "",
		'&':  "No. Diamond",
		'\'': "Crash Site",
		'(':  "Cloudy",
		')':  "MEO",
		'*':  "Snow",
		'+':  "Church",
		',':  "Girl Scout",
		'-':  "Home (HF)",
		'.':  "Unknown Position",
		'/':  "Destination",
		'0':  "No. Circle",
		'1':  "",
		'2':  "",
		'3':  "",
		'4':  "",
		'5':  "",
		'6':  "",
		'7':  "",
		'8':  "",
		'9':  "",
		':':  "Hail",
		';':  "Park",
		'<':  "Gale Flag",
		'=':  "",
		'>':  "No. Car",
		'?':  "Info Kiosk",
		'@':  "Hurricane",
		'A':  "No. Box",
		'B':  "Snow Blowing",
		'C':  "Coast Guard",
		'D':  "Drizzle",
		'E':  "Smoke",
		'F':  "Freeze Rain",
		'G':  "Snow Shower",
		'H':  "Haze",
		'I':  "RainShower",
		'J':  "Lightning",
		'K':  "Kenwood",
		'L':  "Lighthouse",
		'M':  "",
		'N':  "Nav Buoy",
		'O':  "Rocket",
		'P':  "Parking",
	}
)
