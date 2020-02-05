package util

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrZip    = fmt.Errorf("bad zip\n")
	err       = fmt.Errorf("bad Entity")
	regex     = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	LengthErr = fmt.Errorf("bad password, must contain 8-24 characters")
	NumberErr = fmt.Errorf("bad password, must contain at least 1 number")
	CaseErr   = fmt.Errorf("bad password, must contain at least 1 uppercase letter")
	CharErr   = fmt.Errorf("bad password, must contain at least 1 special character")
)

func ValidateEmail(s string) error {
	if !regex.MatchString(s) {
		return err
	}
	return nil
}

// Thanks https://stackoverflow.com/a/25840157.
func ValidatePassword(s string) error {
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

func TypeName(i interface{}) string {
	return strings.Split(TypeOf(i), ".")[0]
}

func TypeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func ValidateZipCode(s string) error {
	if _, err := strconv.Atoi(s); err != nil || len(s) < 5 {
		return ErrZip
	}
	return nil
}

func ZipFromAddressId(s string) string {
	add, _ := base64.StdEncoding.DecodeString(s)
	csv := strings.Split(string(add), ", ")
	return strings.Split(csv[len(csv)-2], "-")[0]
}

func StateFromAddressId(s string) string {
	add, _ := base64.StdEncoding.DecodeString(s)
	csv := strings.Split(string(add), ", ")
	return csv[len(csv)-3]
}
