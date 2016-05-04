package aprs

import (
	"errors"
	"strings"
)

const base91 = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{"

var (
	ErrBase91Decode = errors.New(`aprs: invalid base-91 encoding`)
)

func base91Decode(s string) (n int, err error) {
	for _, c := range s {
		i := strings.IndexRune(base91, c)
		if i < 0 {
			return 0, ErrBase91Decode
		}
		n *= 91
		n += i
	}
	return
}
