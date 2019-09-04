package aprs

import (
	"fmt"
	"time"
)

type TimeFormatError struct {
	Time string
}

func (err TimeFormatError) Error() string {
	return fmt.Sprintf("aprs: unknown time stamp %q", err.Time)
}

func ParseTime(s string) (time.Time, error) {
	if len(s) < 7 {
		return time.Time{}, TimeFormatError{s}
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
	default:
		return time.Time{}, TimeFormatError{s}
	}
}
