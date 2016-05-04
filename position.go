package aprs

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pd0mz/go-maidenhead"
)

var (
	ErrInvalidPosition = errors.New(`aprs: invalid position`)

	// Position ambiguity replacement
	disambiguation = []int{2, 3, 5, 6, 12, 13, 15, 16}
)

const gridChars = "ABCDEFGHIJKLMNOPQRSTUVWX0123456789"

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

func ParseUncompressedPosition(s string) (pos Position, txt string, err error) {
	// APRS PROTOCOL REFERENCE 1.0.1 Chapter 8, page 32 (42 in PDF)

	if len(s) < 18 {
		err = ErrInvalidPosition
		return
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
		latDeg, latMin, latMinFrag uint64
		lngDeg, lngMin, lngMinFrag uint64
		latHemi, lngHemi           byte
		isSouth, isWest            bool
	)

	/* 3210.70N/13132.15E# */
	log.Printf("s: %q\n", s[:18])
	if latDeg, err = strconv.ParseUint(s[0:2], 10, 8); err != nil {
		return
	}
	if latMin, err = strconv.ParseUint(s[2:4], 10, 8); err != nil {
		return
	}
	if latMinFrag, err = strconv.ParseUint(s[5:7], 10, 8); err != nil {
		return
	}
	latHemi = s[7]
	pos.Symbol[0] = s[8]
	if lngDeg, err = strconv.ParseUint(s[9:12], 10, 8); err != nil {
		return
	}
	if lngMin, err = strconv.ParseUint(s[12:14], 10, 8); err != nil {
		return
	}
	if lngMinFrag, err = strconv.ParseUint(s[15:17], 10, 8); err != nil {
		return
	}
	lngHemi = s[17]
	pos.Symbol[1] = s[18]

	if latHemi == 'S' || latHemi == 's' {
		isSouth = true
	} else if latHemi != 'N' && latHemi != 'n' {
		err = ErrInvalidPosition
		return
	}

	if lngHemi == 'W' || lngHemi == 'w' {
		isWest = true
	} else if lngHemi != 'E' && lngHemi != 'e' {
		err = ErrInvalidPosition
		return
	}

	if latDeg > 89 || lngDeg > 179 {
		err = ErrInvalidPosition
		return
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
		txt = s[19:]
	}

	return
}

func ParseCompressedPosition(s string) (pos Position, txt string, err error) {
	// APRS PROTOCOL REFERENCE 1.0.1 Chapter 9, page 36 (46 in PDF)

	if len(s) < 10 {
		err = ErrInvalidPosition
		return
	}

	// Base-91 check
	for _, c := range s[1:9] {
		if c < 0x21 || c > 0x7b {
			err = ErrInvalidPosition
			return
		}
	}

	var lat, lng int
	if lat, err = base91Decode(s[1:5]); err != nil {
		return
	}
	if lng, err = base91Decode(s[5:9]); err != nil {
		return
	}

	pos.Latitude = 90.0 - float64(lat)/380926.0
	pos.Longitude = -180.0 + float64(lng)/190463.0
	pos.Compressed = true
	txt = s[10:]

	return
}

func ParsePositionGrid(s string) (pos Position, txt string, err error) {
	var o int
	for o = 0; o < len(s); o++ {
		if strings.IndexByte(gridChars, s[o]) < 0 {
			break
		}
	}

	if o == 2 || o == 4 || o == 6 || o == 8 {
		var p maidenhead.Point
		if p, err = maidenhead.ParseLocator(s[:o]); err != nil {
			return
		}
		pos.Latitude = p.Latitude
		pos.Longitude = p.Longitude
	}

	if o < len(s) {
		txt = s[o+1:]
	}
	return
}

func ParsePosition(s string, compressed bool) (Position, string, error) {
	log.Printf("parse position %q, %t\n", s, compressed)
	if compressed {
		return ParseCompressedPosition(s)
	}
	return ParseUncompressedPosition(s)
}

func ParsePositionBoth(s string) (pos Position, txt string, err error) {
	if pos, txt, err = ParseUncompressedPosition(s); err != nil {
		pos, txt, err = ParseCompressedPosition(s)
	}
	return
}
