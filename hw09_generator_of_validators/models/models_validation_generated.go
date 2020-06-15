// Code generated automatically. DO NOT EDIT.
package models

import (
	"errors"
	"regexp"
)

var (
	ErrMin    = errors.New("Value is less than min")
	ErrMax    = errors.New("Value is greater than max")
	ErrLen    = errors.New("Length validation error")
	ErrRegexp = errors.New("Regexp validation error")
	ErrEnum   = errors.New("Enum validation error")
)

type ValidationError struct {
	Field string
	Err   error
}

func (s App) Validate() (verr []ValidationError, err error) {

	version := s.Version

	if len(version) != 5 {
		ve := ValidationError{
			Field: "Version",
			Err:   ErrLen,
		}
		verr = append(verr, ve)

	}

	return
}

func (s Response) Validate() (verr []ValidationError, err error) {

	code := s.Code

	if code != 200 && code != 404 && code != 500 {
		ve := ValidationError{
			Field: "Code",
			Err:   ErrEnum,
		}
		verr = append(verr, ve)

	}

	return
}

func (s User) Validate() (verr []ValidationError, err error) {

	id := s.ID

	if len(id) != 36 {
		ve := ValidationError{
			Field: "ID",
			Err:   ErrLen,
		}
		verr = append(verr, ve)

	}

	age := s.Age

	if age < 18 {
		ve := ValidationError{
			Field: "Age",
			Err:   ErrMin,
		}
		verr = append(verr, ve)

	}

	if age > 50 {
		ve := ValidationError{
			Field: "Age",
			Err:   ErrMax,
		}
		verr = append(verr, ve)

	}

	var re = regexp.MustCompile(`^\w+@\w+\.\w+$`)

	email := s.Email

	if !re.MatchString(email) {
		ve := ValidationError{
			Field: "Email",
			Err:   ErrRegexp,
		}
		verr = append(verr, ve)

	}

	role := s.Role

	if role != "admin" && role != "stuff" {
		ve := ValidationError{
			Field: "Role",
			Err:   ErrEnum,
		}
		verr = append(verr, ve)

	}

	for _, phones := range s.Phones {

		if len(phones) != 11 {
			ve := ValidationError{
				Field: "Phones",
				Err:   ErrLen,
			}
			verr = append(verr, ve)

			break

		}

	}

	return
}
