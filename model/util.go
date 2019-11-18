package model

import (
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"unicode"
)

// No email regex is perfect, but this one is close.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Determines email address validity against predefined regex.
// Allows addresses with third party domains and any extension.
func IsEmailValid(s string) error {
	if emailRegex.MatchString(s) == false {
		return fmt.Errorf("bad email [%s]", s)
	} else {
		return nil
	}
}

// The following is an adaptation of https://stackoverflow.com/a/25840157
func IsPasswordValid(s string) error {
	var number, upper, special bool
	length := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			length++
		case unicode.IsUpper(c):
			upper = true
			length++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			length++
		case unicode.IsLetter(c) || c == ' ':
			length++
		default:
			// do not increment length for unrecognized characters
		}
	}
	if length < 7 || length > 24 {
		return fmt.Errorf("bad password, must contain 8-24 characters")
	} else if number == false {
		return fmt.Errorf("bad password, must contain at least 1 number")
	} else if upper == false {
		return fmt.Errorf("bad password, must contain at least 1 uppercase letter")
	} else if special == false {
		return fmt.Errorf("bad password, must contain at least 1 special character")
	} else {
		return nil
	}
}

func IsIdValid(s string) error {
	if _, err := uuid.Parse(s); err != nil {
		return err
	} else {
		return nil
	}
}
