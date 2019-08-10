package aprs

import (
	"errors"
	"strings"
)

const base91 = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{"

func base91Decode(s string) (int, error) {
	var n int
	for _, c := range s {
		i := strings.IndexRune(base91, c)
		if i < 0 {
			return 0, errors.New("aprs: invalid base-91 encoding")
		}
		n *= 91
		n += i
	}
	return n, nil
}
