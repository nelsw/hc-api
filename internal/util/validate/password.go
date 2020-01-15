package validate

import (
	"fmt"
	"unicode"
)

var (
	LengthErr = fmt.Errorf("bad password, must contain 8-24 characters")
	NumberErr = fmt.Errorf("bad password, must contain at least 1 number")
	CaseErr   = fmt.Errorf("bad password, must contain at least 1 uppercase letter")
	CharErr   = fmt.Errorf("bad password, must contain at least 1 special character")
)

func Password(s string) error {
	// Does the password meet predefined criteria?
	// Thanks https://stackoverflow.com/a/25840157.
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
		}
	}
	if length < 8 || length > 24 {
		return LengthErr
	} else if !number {
		return NumberErr
	} else if !upper {
		return CaseErr
	} else if !special {
		return CharErr
	} else {
		return nil
	}
}
