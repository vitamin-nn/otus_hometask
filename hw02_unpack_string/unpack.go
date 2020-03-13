package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	var prev rune

	for i, curr := range s {
		if unicode.IsDigit(curr) {
			if unicode.IsDigit(prev) || i == 0 {
				return "", ErrInvalidString
			}
			//rep, _ := strconv.Atoi(string(curr))
			rep := int(curr - '0')
			w := strings.Repeat(string(prev), (rep - 1))
			b.WriteString(w)
		} else {
			b.WriteRune(curr)
		}
		prev = curr
	}
	result := b.String()
	return result, nil
}
