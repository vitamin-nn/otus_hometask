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
		var w string   // will written in result
		if isEscaped { // if backslashed current symbol
			w = string(curr)
			isEscaped = false
			isDigitAllowed = true
		} else if unicode.IsDigit(curr) {
			factor := int(curr - '0')
			// check errs: two digit in a row, first sign is digit, 0 sign
			if !isDigitAllowed || i == 0 || factor == 0 {
				return "", ErrInvalidString
			}
			w = strings.Repeat(string(prev), (factor - 1))
			isDigitAllowed = false
		} else if string(curr) == "\\" {
			isDigitAllowed = true
			isEscaped = true
		} else {
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

// UnpackSimple without processing of backslash
func UnpackSimple(s string) (string, error) {
	var b strings.Builder
	var prev rune

	for i, curr := range s {
		var w string
		if unicode.IsDigit(curr) {
			rep := int(curr - '0')
			// check errs: two digit in a row, first sign is digit, 0 sign
			if unicode.IsDigit(prev) || i == 0 || rep == 0 {
				return "", ErrInvalidString
			}
			w = strings.Repeat(string(prev), (rep - 1))
		} else {
			w = string(curr)
		}
		if w != "" {
			b.WriteString(w)
		}
		prev = curr
	}
	result := b.String()
	return result, nil
}
