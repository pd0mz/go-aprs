package aprs

import "fmt"

type Symbol [2]byte

func (s Symbol) IsPrimaryTable() bool { return s[0] != '\\' }

func (s Symbol) get(idx int) (string, error) {
	var m map[byte]map[int]string
	if s.IsPrimaryTable() {
		m = primarySymbol
	} else {
		m = alternateSymbol
	}
	n, ok := m[s[1]]
	if !ok {
		return "", fmt.Errorf("unknown symbol %c", s[1])
	}
	if i, ok := n[idx]; ok {
		return i, nil
	}
	return "", fmt.Errorf("symbol doesn't have requested index: %v", n)
}

func (s Symbol) String() string {
	hr, err := s.get(1)
	if err != nil {
		return err.Error()
	}
	return hr
}

func (s Symbol) SSID() (string, error) {
	return s.get(2)
}

func (s Symbol) Emoji() (string, error) {
	return s.get(3)
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

var (
	// Source: http://www.aprs.org/symbols/symbolsX.txt
	// 0: XYZ code
	// 1: Human readable
	// 2: SSID
	// 3: Emoji
	primarySymbol = map[byte]map[int]string{
		'!':  map[int]string{0: "BB", 1: "Police, Sheriff", 3: ":cop:"},
		'"':  map[int]string{0: "BC", 1: "reserved"},
		'#':  map[int]string{0: "BD", 1: "Digi"},
		'$':  map[int]string{0: "BE", 1: "Phone", 3: ":phone:"},
		'%':  map[int]string{0: "BF", 1: "DX Cluster"},
		'&':  map[int]string{0: "BG", 1: "HF Gateway"},
		'\'': map[int]string{0: "BH", 1: "Small Aircraft", 2: "11", 3: ":airplane:"},
		'(':  map[int]string{0: "BI", 1: "Mobile Satellite Station", 3: ":satellite:"},
		')':  map[int]string{0: "BJ", 1: "Wheelchair", 3: ":wheelchair:"},
		'*':  map[int]string{0: "BK", 1: "Snowmobile"},
		'+':  map[int]string{0: "BL", 1: "Red Cross"},
		',':  map[int]string{0: "BM", 1: "Boy Scout"},
		'-':  map[int]string{0: "BN", 1: "House QTH (VHF)"},
		'.':  map[int]string{0: "BO", 1: "X"},
		'/':  map[int]string{0: "BP", 1: "Red Dot"},
		'0':  map[int]string{0: "P0", 1: "Circle (0)"},
		'1':  map[int]string{0: "P1", 1: "Circle (1)"},
		'2':  map[int]string{0: "P2", 1: "Circle (2)"},
		'3':  map[int]string{0: "P3", 1: "Circle (3)"},
		'4':  map[int]string{0: "P4", 1: "Circle (4)"},
		'5':  map[int]string{0: "P5", 1: "Circle (5)"},
		'6':  map[int]string{0: "P6", 1: "Circle (6)"},
		'7':  map[int]string{0: "P7", 1: "Circle (7)"},
		'8':  map[int]string{0: "P8", 1: "Circle (8)"},
		'9':  map[int]string{0: "P9", 1: "Circle (9)"},
		':':  map[int]string{0: "MR", 1: "Fire", 3: ":fire:"},
		';':  map[int]string{0: "MS", 1: "Campground", 3: ":tent:"},
		'<':  map[int]string{0: "MT", 1: "Motorcycle", 2: "10", 3: ":bike:"},
		'=':  map[int]string{0: "MU", 1: "Railroad Engine", 3: ":train:"},
		'>':  map[int]string{0: "MV", 1: "Car", 2: "9", 3: ":car:"},
		'?':  map[int]string{0: "MW", 1: "File Server"},
		'@':  map[int]string{0: "MX", 1: "HC Future"},
		'A':  map[int]string{0: "PA", 1: "Aid Station", 3: ":hospital:"},
		'B':  map[int]string{0: "PB", 1: "BBS or PBBS"},
		'C':  map[int]string{0: "PC", 1: "Canoe"},
		'D':  map[int]string{0: "PD"},
		'E':  map[int]string{0: "PE", 1: "Eyeball"},
		'F':  map[int]string{0: "PF", 1: "Tractor", 3: ":tractor:"},
		'G':  map[int]string{0: "PG", 1: "Grid Square"},
		'H':  map[int]string{0: "PH", 1: "Hotel", 3: ":hotel:"},
		'I':  map[int]string{0: "PI", 1: "TCP/IP"},
		'J':  map[int]string{0: "PJ"},
		'K':  map[int]string{0: "PK", 1: "School", 3: ":school:"},
		'L':  map[int]string{0: "PL", 1: "PC User", 3: ":computer:"},
		'M':  map[int]string{0: "PM", 1: "MacAPRS", 3: ":computer:"},
		'N':  map[int]string{0: "PN", 1: "NTS Station"},
		'O':  map[int]string{0: "PO", 1: "Balloon", 2: "11", 3: ":airplane:"},
		'P':  map[int]string{0: "PP", 1: "Police", 3: ":police_car:"},
		'Q':  map[int]string{0: "PQ"},
		'R':  map[int]string{0: "PR", 1: "Recreational Vehicle", 2: "13", 3: ":car:"},
		'S':  map[int]string{0: "PS", 1: "Shuttle"},
		'T':  map[int]string{0: "PT", 1: "SSTV"},
		'U':  map[int]string{0: "PU", 1: "Bus", 2: "2", 3: ":bus:"},
		'V':  map[int]string{0: "PV", 1: "ATV"},
		'W':  map[int]string{0: "PW", 1: "National WX Service Site"},
		'X':  map[int]string{0: "PX", 1: "Helo", 2: "6"},
		'Y':  map[int]string{0: "PY", 1: "Yacht", 2: "5", 3: ":sailboat:"},
		'Z':  map[int]string{0: "PZ", 1: "WinAPRS", 3: ":computer:"},
		'[':  map[int]string{0: "HS", 1: "Human/Person", 2: "7", 3: ":running:"},
		'\\': map[int]string{0: "HT", 1: "DF Station"},
		']':  map[int]string{0: "HU", 1: "Post Office", 3: ":post_office:"},
		'^':  map[int]string{0: "HV", 1: "Large Aircraft", 3: ":airplane:"},
		'_':  map[int]string{0: "HW", 1: "Weather Station", 3: ":cloud:"},
		'`':  map[int]string{0: "HX", 1: "Dish Antenna", 3: ":satellite:"},
		'a':  map[int]string{0: "LA", 1: "Ambulance", 2: "1", 3: ":ambulance:"},
		'b':  map[int]string{0: "LB", 1: "Bike", 2: "4", 3: ":bike:"},
		'c':  map[int]string{0: "LC", 1: "Incident Command Post"},
		'd':  map[int]string{0: "LD", 1: "Fire Dept", 3: ":fire_engine:"},
		'e':  map[int]string{0: "LE", 1: "Horse", 3: ":racehorse:"},
		'f':  map[int]string{0: "LF", 1: "Fire Truck", 2: "3", 3: ":fire_engine:"},
		'g':  map[int]string{0: "LG", 1: "Glider", 3: ":airplane:"},
		'h':  map[int]string{0: "LH", 1: "Hospital", 3: ":hospital:"},
		'i':  map[int]string{0: "LI", 1: "IOTA"},
		'j':  map[int]string{0: "LJ", 1: "Jeep", 2: "12", 3: ":car:"},
		'k':  map[int]string{0: "LK", 1: "Truck", 2: "14", 3: ":truck:"},
		'l':  map[int]string{0: "LL", 1: "Laptop", 3: ":computer:"},
		'm':  map[int]string{0: "LM", 1: "Mic-E Repeater"},
		'n':  map[int]string{0: "LN", 1: "Node"},
		'o':  map[int]string{0: "LO", 1: "EOC"},
		'p':  map[int]string{0: "LP", 1: "Dog", 3: ":dog2:"},
		'q':  map[int]string{0: "LQ", 1: "Grid SQ"},
		'r':  map[int]string{0: "LR", 1: "Repeater"},
		's':  map[int]string{0: "LS", 1: "Ship", 2: "8", 3: ":ship:"},
		't':  map[int]string{0: "LT", 1: "Truck Stop"},
		'u':  map[int]string{0: "LU", 1: "Truck (18 Wheeler)", 3: ":truck:"},
		'v':  map[int]string{0: "LV", 1: "Van", 2: "15", 3: ":minibus:"},
		'w':  map[int]string{0: "LW", 1: "Water Station"},
		'x':  map[int]string{0: "LX", 1: "xAPRS", 3: ":computer:"},
		'y':  map[int]string{0: "LY", 1: "Yagi @ QTH"},
		'z':  map[int]string{0: "LZ"},
		'{':  map[int]string{0: "J1"},
		'|':  map[int]string{0: "J2", 1: "TNC Stream Switch"},
		'}':  map[int]string{0: "J3"},
		'~':  map[int]string{0: "J4", 1: "TNC Stream Switch"},
	}
	alternateSymbol = map[byte]map[int]string{
		'!':  map[int]string{0: "OBO", 1: "Emergency"},
		'"':  map[int]string{0: "OC", 1: "Reserved"},
		'#':  map[int]string{0: "OD#", 1: "Overlay Digi"},
		'$':  map[int]string{0: "OEO", 1: "Bank/ATM", 3: ":atm:"},
		'%':  map[int]string{0: "OFO", 1: "Power Plant", 3: ":factory:"},
		'&':  map[int]string{0: "OG#", 1: "I=Igte R=RX T=1hopTX 2=2hopTX"},
		'\'': map[int]string{0: "OHO", 1: "Crash Site"},
		'(':  map[int]string{0: "OIO", 1: "Cloudy", 3: ":cloud:"},
		')':  map[int]string{0: "OJO", 1: "Firenet MEO"},
		'*':  map[int]string{0: "OK"},
		'+':  map[int]string{0: "OL", 1: "Church", 3: ":church:"},
		',':  map[int]string{0: "OM", 1: "Girl Scouts", 3: ":tent:"},
		'-':  map[int]string{0: "ONO", 1: "House", 3: ":house:"},
		'.':  map[int]string{0: "OO", 1: "Ambiguous"},
		'/':  map[int]string{0: "OP", 1: "Waypoint Destination"},
		'0':  map[int]string{0: "A0#", 1: "Circle", 3: ":red_circle:"},
		'1':  map[int]string{0: "A1"},
		'2':  map[int]string{0: "A2"},
		'3':  map[int]string{0: "A3"},
		'4':  map[int]string{0: "A4"},
		'5':  map[int]string{0: "A5"},
		'6':  map[int]string{0: "A6"},
		'7':  map[int]string{0: "A7"},
		'8':  map[int]string{0: "A8O", 1: "WiFi Network"},
		'9':  map[int]string{0: "A9", 1: "Gas Station", 3: ":fuelpump:"},
		':':  map[int]string{0: "NR"},
		';':  map[int]string{0: "NSO", 1: "Park/Picnic"},
		'<':  map[int]string{0: "NTO", 1: "Advisory"},
		'=':  map[int]string{0: "NUO"},
		'>':  map[int]string{0: "NV#", 1: "Cars & Vehicles", 3: ":car:"},
		'?':  map[int]string{0: "NW", 1: "Info Kiosk"},
		'@':  map[int]string{0: "NX", 1: "Hurricane", 3: ":cyclone:"},
		'A':  map[int]string{0: "AA#", 1: "Box DTMF & RFID"},
		'B':  map[int]string{0: "AB"},
		'C':  map[int]string{0: "AC", 1: "Coast Guard"},
		'D':  map[int]string{0: "ADO", 1: "Depots"},
		'E':  map[int]string{0: "AE", 1: "Smoke"},
		'F':  map[int]string{0: "AF"},
		'G':  map[int]string{0: "AG"},
		'H':  map[int]string{0: "AHO", 1: "Haze"},
		'I':  map[int]string{0: "AI", 1: "Rain Shower", 3: ":umbrella:"},
		'J':  map[int]string{0: "AJ"},
		'K':  map[int]string{0: "AK", 1: "Kenwood HT"},
		'L':  map[int]string{0: "AL", 1: "Lighthouse"},
		'M':  map[int]string{0: "AMO", 1: "MARS"},
		'N':  map[int]string{0: "AN", 1: "Navigation Buoy"},
		'O':  map[int]string{0: "AO", 1: "Rocket", 3: ":rocket:"},
		'P':  map[int]string{0: "AP", 1: "Parking", 3: ":parking:"},
		'Q':  map[int]string{0: "AQ", 1: "Quake"},
		'R':  map[int]string{0: "ARO", 1: "Restaurant"},
		'S':  map[int]string{0: "AS", 1: "Satellite/Pacsat", 3: ":rocket:"},
		'T':  map[int]string{0: "AT", 1: "Thunderstorm"},
		'U':  map[int]string{0: "AU", 1: "Sunny"},
		'V':  map[int]string{0: "AV", 1: "VORTAC Nav Aid"},
		'W':  map[int]string{0: "AW#", 1: "NWS Site"},
		'X':  map[int]string{0: "AX", 1: "Pharmacy"},
		'Y':  map[int]string{0: "AYO", 1: "Radios and devices"},
		'Z':  map[int]string{0: "AZ"},
		'[':  map[int]string{0: "DSO", 1: "W. Cloud"},
		'\\': map[int]string{0: "DTO", 1: "GPS"},
		']':  map[int]string{0: "DU"},
		'^':  map[int]string{0: "DV#", 1: "Other Aircraft", 3: ":airplane:"},
		'_':  map[int]string{0: "DW#", 1: "WX Site"},
		'`':  map[int]string{0: "DX", 1: "Rain", 3: ":umbrella:"},
		'a':  map[int]string{0: "SA#O"},
		'b':  map[int]string{0: "SB"},
		'c':  map[int]string{0: "SC#O", 1: "CD Triangle"},
		'd':  map[int]string{0: "SD", 1: "DX Spot"},
		'e':  map[int]string{0: "SE", 1: "Sleet"},
		'f':  map[int]string{0: "SF", 1: "Funnel Cloud"},
		'g':  map[int]string{0: "SG", 1: "Gale Flags"},
		'h':  map[int]string{0: "SHO", 1: "Store or Hamfest"},
		'i':  map[int]string{0: "SI#", 1: "Box / POI"},
		'j':  map[int]string{0: "SJ", 1: "Work Zone"},
		'k':  map[int]string{0: "SKO", 1: "Special Vehicle"},
		'l':  map[int]string{0: "SL", 1: "Areas"},
		'm':  map[int]string{0: "SM", 1: "Value Sign"},
		'n':  map[int]string{0: "SN#", 1: "Triangle"},
		'o':  map[int]string{0: "SO", 1: "Small Circle"},
		'p':  map[int]string{0: "SP"},
		'q':  map[int]string{0: "SQ"},
		'r':  map[int]string{0: "SR", 1: "Restrooms"},
		's':  map[int]string{0: "SS#", 1: "Ship/Boats", 3: ":speedboat:"},
		't':  map[int]string{0: "ST", 1: "Tornado", 3: ":cyclone:"},
		'u':  map[int]string{0: "SU#", 1: "Truck", 3: ":truck:"},
		'v':  map[int]string{0: "SV#", 1: "Van", 3: ":minibus:"},
		'w':  map[int]string{0: "SWO", 1: "Flooding"},
		'x':  map[int]string{0: "SX", 1: "Wreck/Obstruction"},
		'y':  map[int]string{0: "SY", 1: "Skywarn"},
		'z':  map[int]string{0: "SZ#", 1: "Shelter"},
		'{':  map[int]string{0: "Q1"},
		'|':  map[int]string{0: "Q2", 1: "TNC Stream Switch"},
		'}':  map[int]string{0: "Q3"},
		'~':  map[int]string{0: "Q4", 1: "TNC Stream Switch"},
	}
)
