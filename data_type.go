package aprs

import "fmt"

type DataType byte

var (
	dataTypeName = map[DataType]string{
		0x1c: "Current Mic-E Data (Rev 0 beta)",
		0x1d: "Old Mic-E Data (Rev 0 beta)",
		'!':  "Position without timestamp (no APRS messaging), or Ultimeter 2000 WX Station",
		'#':  "Peet Bros U-II Weather Station",
		'$':  "Raw GPS data or Ultimeter 2000",
		'%':  "Agrelo DFJr / MicroFinder",
		'"':  "Old Mic-E Data (but Current data for TM-D700)",
		')':  "Item",
		'*':  "Peet Bros U-II Weather Station",
		',':  "Invalid data or test data",
		'/':  "Position with timestamp (no APRS messaging)",
		':':  "Message",
		';':  "Object",
		'<':  "Station Capabilities",
		'=':  "Position without timestamp (with APRS messaging)",
		'>':  "Status",
		'?':  "Query",
		'@':  "Position with timestamp (with APRS messaging)",
		'T':  "Telemetry data",
		'[':  "Maidenhead grid locator beacon (obsolete)",
		'_':  "Weather Report (without position)",
		'`':  "Current Mic-E Data (not used in TM-D700)",
		'{':  "User-Defined APRS packet format",
		'}':  "Third-party traffic",
	}
)

func (t DataType) String() string {
	if s, ok := dataTypeName[t]; ok {
		return s
	}
	return fmt.Sprintf("Unknown packet type %#02x", byte(t))
}
