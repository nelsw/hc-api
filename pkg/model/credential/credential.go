package credential

import (
	"fmt"
	"regexp"
	"unicode"
)

type Entity struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	UserId   string `json:"user_id"`
}

var pattern = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
var regex = regexp.MustCompile(pattern)

func (e *Entity) Validate() error {
	if err := validateEmail(e.Id); err != nil {
		return err
	} else if err := validatePassword(e.Password); err != nil {
		return err
	} else {
		return nil
	}
}

func validateEmail(s string) error {
	if !regex.MatchString(s) {
		return fmt.Errorf("bad email, did not match pattern [%s]", pattern)
	}
	return nil
}

// Thanks https://stackoverflow.com/a/25840157.
func validatePassword(s string) error {
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
		return fmt.Errorf("bad password, must contain 8-24 characters")
	} else if !number {
		return fmt.Errorf("bad password, must contain at least 1 number")
	} else if !upper {
		return fmt.Errorf("bad password, must contain at least 1 uppercase letter")
	} else if !special {
		return fmt.Errorf("bad password, must contain at least 1 special character")
	} else {
		return nil
	}
}
