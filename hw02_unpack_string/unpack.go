package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// with processing of backslash
func Unpack(s string) (string, error) {
	var b strings.Builder
	var prev rune
	isDigitAllowed := true
	isEscaped := false

	for i, curr := range s {
		var w string // will written in result
		switch {
		case isEscaped:
			w = string(curr)
			isEscaped = false
			isDigitAllowed = true
		case unicode.IsDigit(curr):
			factor := int(curr - '0')
			// check errs: two digit in a row, first sign is digit, 0 sign
			if !isDigitAllowed || i == 0 || factor == 0 {
				return "", ErrInvalidString
			}
			w = strings.Repeat(string(prev), (factor - 1))
			isDigitAllowed = false
		case string(curr) == "\\":
			isDigitAllowed = true
			isEscaped = true
		default:
			w = string(curr)
			isDigitAllowed = true
		}
		if w != "" {
			b.WriteString(w)
		}
		prev = curr
	}
	result := b.String()
	return result, nil
}
