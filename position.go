package aprs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pd0mz/go-maidenhead"
)

var (
	// Position ambiguity replacement
	disambiguation = []int{2, 3, 5, 6, 12, 13, 15, 16}

	miceCodes = map[rune]map[int]string{
		'0': map[int]string{0: "0", 1: "0", 2: "S", 3: "0", 4: "E"},
		'1': map[int]string{0: "1", 1: "0", 2: "S", 3: "0", 4: "E"},
		'2': map[int]string{0: "2", 1: "0", 2: "S", 3: "0", 4: "E"},
		'3': map[int]string{0: "3", 1: "0", 2: "S", 3: "0", 4: "E"},
		'4': map[int]string{0: "4", 1: "0", 2: "S", 3: "0", 4: "E"},
		'5': map[int]string{0: "5", 1: "0", 2: "S", 3: "0", 4: "E"},
		'6': map[int]string{0: "6", 1: "0", 2: "S", 3: "0", 4: "E"},
		'7': map[int]string{0: "7", 1: "0", 2: "S", 3: "0", 4: "E"},
		'8': map[int]string{0: "8", 1: "0", 2: "S", 3: "0", 4: "E"},
		'9': map[int]string{0: "9", 1: "0", 2: "S", 3: "0", 4: "E"},
		'A': map[int]string{0: "0", 1: "1 (Custom)"},
		'B': map[int]string{0: "1", 1: "1 (Custom)"},
		'C': map[int]string{0: "2", 1: "1 (Custom)"},
		'D': map[int]string{0: "3", 1: "1 (Custom)"},
		'E': map[int]string{0: "4", 1: "1 (Custom)"},
		'F': map[int]string{0: "5", 1: "1 (Custom)"},
		'G': map[int]string{0: "6", 1: "1 (Custom)"},
		'H': map[int]string{0: "7", 1: "1 (Custom)"},
		'I': map[int]string{0: "8", 1: "1 (Custom)"},
		'J': map[int]string{0: "9", 1: "1 (Custom)"},
		'K': map[int]string{0: " ", 1: "1 (Custom)"},
		'L': map[int]string{0: " ", 1: "0", 2: "S", 3: "0", 4: "E"},
		'P': map[int]string{0: "0", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'Q': map[int]string{0: "1", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'R': map[int]string{0: "2", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'S': map[int]string{0: "3", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'T': map[int]string{0: "4", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'U': map[int]string{0: "5", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'V': map[int]string{0: "6", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'W': map[int]string{0: "7", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'X': map[int]string{0: "8", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'Y': map[int]string{0: "9", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
		'Z': map[int]string{0: " ", 1: "1 (Std)", 2: "N", 3: "100", 4: "W"},
	}

	miceMsgTypes = map[string]string{
		"000":          "Emergency",
		"001 (Std)":    "Priority",
		"001 (Custom)": "Custom-6",
		"010 (Std)":    "Special",
		"010 (Custom)": "Custom-5",
		"011 (Std)":    "Committed",
		"011 (Custom)": "Custom-4",
		"100 (Std)":    "Returning",
		"100 (Custom)": "Custom-3",
		"101 (Std)":    "In Service",
		"101 (Custom)": "Custom-2",
		"110 (Std)":    "En Route",
		"110 (Custom)": "Custom-1",
		"111 (Std)":    "Off Duty",
		"111 (Custom)": "Custom-0",
	}
)

const (
	gridChars = "ABCDEFGHIJKLMNOPQRSTUVWX0123456789"

	messageTypeStd    = "Std"
	messageTypeCustom = "Custom"
)

type Position struct {
	Latitude   float64 // Degrees
	Longitude  float64 // Degrees
	Ambiguity  int
	Symbol     Symbol
	Compressed bool
}

func (pos Position) String() string {
	if pos.Ambiguity == 0 {
		return fmt.Sprintf("{%f, %f}", pos.Latitude, pos.Longitude)
	}
	return fmt.Sprintf("{%f, %f}, ambiguity=%d", pos.Latitude, pos.Longitude, pos.Ambiguity)
}

func ParseUncompressedPosition(s string) (Position, string, error) {
	// APRS PROTOCOL REFERENCE 1.0.1 Chapter 8, page 32 (42 in PDF)

	pos := Position{}

	if len(s) < 18 {
		return pos, "", errors.New("aprs: invalid position")
	}

	b := []byte(s)
	for _, p := range disambiguation {
		if b[p] == ' ' {
			pos.Ambiguity++
			b[p] = '0'
		}
	}
	s = string(b)

	var (
		err                        error
		latDeg, latMin, latMinFrag uint64
		lngDeg, lngMin, lngMinFrag uint64
		latHemi, lngHemi           byte
		isSouth, isWest            bool
	)

	if latDeg, err = strconv.ParseUint(s[0:2], 10, 8); err != nil {
		return pos, "", err
	}
	if latMin, err = strconv.ParseUint(s[2:4], 10, 8); err != nil {
		return pos, "", err
	}
	if latMinFrag, err = strconv.ParseUint(s[5:7], 10, 8); err != nil {
		return pos, "", err
	}
	latHemi = s[7]
	pos.Symbol[0] = s[8]
	if lngDeg, err = strconv.ParseUint(s[9:12], 10, 8); err != nil {
		return pos, "", err
	}
	if lngMin, err = strconv.ParseUint(s[12:14], 10, 8); err != nil {
		return pos, "", err
	}
	if lngMinFrag, err = strconv.ParseUint(s[15:17], 10, 8); err != nil {
		return pos, "", err
	}
	lngHemi = s[17]
	pos.Symbol[1] = s[18]

	if latHemi == 'S' || latHemi == 's' {
		isSouth = true
	} else if latHemi != 'N' && latHemi != 'n' {
		return pos, "", errors.New("aprs: invalid position")
	}

	if lngHemi == 'W' || lngHemi == 'w' {
		isWest = true
	} else if lngHemi != 'E' && lngHemi != 'e' {
		return pos, "", errors.New("aprs: invalid position")
	}

	if latDeg > 89 || lngDeg > 179 {
		return pos, "", errors.New("aprs: invalid position")
	}

	pos.Latitude = float64(latDeg) + float64(latMin)/60.0 + float64(latMinFrag)/6000.0
	pos.Longitude = float64(lngDeg) + float64(lngMin)/60.0 + float64(lngMinFrag)/6000.0

	if isSouth {
		pos.Latitude = 0.0 - pos.Latitude
	}
	if isWest {
		pos.Longitude = 0.0 - pos.Longitude
	}

	if pos.Symbol[1] >= 'a' || pos.Symbol[1] <= 'k' {
		pos.Symbol[1] -= 32
	}

	if len(s) > 19 {
		return pos, s[19:], nil
	}
	return pos, "", nil
}

func ParseCompressedPosition(s string) (Position, string, error) {
	// APRS PROTOCOL REFERENCE 1.0.1 Chapter 9, page 36 (46 in PDF)

	pos := Position{}

	if len(s) < 10 {
		return pos, "", errors.New("aprs: invalid position")
	}

	// Base-91 check
	for _, c := range s[1:9] {
		if c < 0x21 || c > 0x7b {
			return pos, "", errors.New("aprs: invalid position")
		}
	}

	var err error
	var lat, lng int
	if lat, err = base91Decode(s[1:5]); err != nil {
		return pos, "", err
	}
	if lng, err = base91Decode(s[5:9]); err != nil {
		return pos, "", err
	}

	pos.Latitude = 90.0 - float64(lat)/380926.0
	pos.Longitude = -180.0 + float64(lng)/190463.0
	pos.Compressed = true

	return pos, s[10:], nil
}

func ParseMicE(s, dest string) (Position, error) {
	// APRS PROTOCOL REFERENCE 1.0.1 Chapter 10, page 42 in PDF

	pos := Position{}

	if len(s) < 9 || len(dest) != 6 {
		return pos, errors.New("aprs: invalid MicE position")
	}

	ns := miceCodes[rune(dest[3])][2]
	we := miceCodes[rune(dest[5])][4]

	latF := fmt.Sprintf("%s%s", miceCodes[rune(dest[0])][0], miceCodes[rune(dest[1])][0])
	latF = strings.Trim(latF, ". ")
	latD, err := strconv.ParseFloat(latF, 64)
	if err != nil {
		return pos, errors.New("aprs: invalid position")
	}
	lonF := fmt.Sprintf("%s%s.%s%s", miceCodes[rune(dest[2])][0], miceCodes[rune(dest[3])][0], miceCodes[rune(dest[4])][0], miceCodes[rune(dest[5])][0])
	lonF = strings.Trim(lonF, ". ")
	latM, err := strconv.ParseFloat(lonF, 64)
	if err != nil {
		return pos, errors.New("aprs: invalid position")
	}
	if latM != 0 {
		latD += latM / 60
	}
	if strings.ToUpper(ns) == "S" {
		latD = -latD
	}

	lonOff := miceCodes[rune(dest[4])][3]
	lonD := float64(s[1]) - 28
	if lonOff == "100" {
		lonD += 100
	}
	if lonD >= 180 && lonD < 190 {
		lonD -= 80
	}
	if lonD >= 190 && lonD < 200 {
		lonD -= 190
	}

	lonM := float64(s[2]) - 28
	if lonM >= 60 {
		lonM -= 60
	}
	// adding hundreth of minute then add minute as deg fraction
	lonH := float64(s[3]) - 28
	if lonH != 0 {
		lonM += lonH / 100
	}
	if lonM != 0 {
		lonD += lonM / 60
	}
	if strings.ToUpper(we) == "W" {
		lonD = -lonD
	}

	pos.Latitude = latD
	pos.Longitude = lonD

	return pos, nil
}

func ParsePositionGrid(s string) (Position, string, error) {
	var o int
	for o = 0; o < len(s); o++ {
		if strings.IndexByte(gridChars, s[o]) < 0 {
			break
		}
	}

	pos := Position{}
	if o == 2 || o == 4 || o == 6 || o == 8 {
		p, err := maidenhead.ParseLocator(s[:o])
		if err != nil {
			return pos, "", err
		}
		pos.Latitude = p.Latitude
		pos.Longitude = p.Longitude
	}

	var txt string
	if o < len(s) {
		txt = s[o+1:]
	}
	return pos, txt, nil
}

func ParsePosition(s string, compressed bool) (Position, string, error) {
	if compressed {
		return ParseCompressedPosition(s)
	}
	return ParseUncompressedPosition(s)
}

func ParsePositionBoth(s string) (Position, string, error) {
	pos, txt, err := ParseUncompressedPosition(s)
	if err != nil {
		return ParseCompressedPosition(s)
	}
	return pos, txt, err
}
