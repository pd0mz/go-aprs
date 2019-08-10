package aprs

import (
	"errors"
	"fmt"
	"time"
)

func ParseTime(s string) (time.Time, error) {
	if len(s) < 7 {
		return time.Time{}, fmt.Errorf("aprs: invalid time: %q is less than 7 chars", s)
	}

	switch {
	case s[6] == 'z': // Day/Hours/Minutes (DHM) format
		return time.Parse("021504", s[:6])
	case s[6] == '/': // Day/Hours/Minutes (DHM) format
		return time.Parse("021504", s[:6])
	case s[6] == 'h': // Hours/Minutes/Seconds (HMS) format
		return time.Parse("150405", s[:6])
	case len(s) >= 8: // Month/Day/Hours/Minutes (MDHM) format
		return time.Parse("01021504", s[:8])
	}

	return time.Time{}, errors.New("aprs: invalid time (unclear format)")
}
