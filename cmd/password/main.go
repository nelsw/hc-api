package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
	. "hc-api/service"
	"os"
	"unicode"
)

var (
	table        = os.Getenv("USER_PASSWORD_TABLE")
	lengthError  = fmt.Errorf("bad password, must contain 8-24 characters")
	numberError  = fmt.Errorf("bad password, must contain at least 1 number")
	upperError   = fmt.Errorf("bad password, must contain at least 1 uppercase letter")
	specialError = fmt.Errorf("bad password, must contain at least 1 special character")
)

// Data structure for persisting and retrieving a users (encrypted) password. The user entity maintains a 1-1 FK
// relationship with the UserPassword entity for referential integrity, where user.PasswordId == userPassword.Id.
type Password struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

// Validates the UserPassword entity by confirming that both the password and id values are valid.
// Most of the following method body is an adaptation of https://stackoverflow.com/a/25840157.
func (p *Password) Validate() error {
	var number, upper, special bool
	length := 0
	for _, c := range p.Password {
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
		return lengthError
	} else if number == false {
		return numberError
	} else if upper == false {
		return upperError
	} else if special == false {
		return specialError
	} else {
		return nil
	}
}

func Handle(p Password) error {

	if err := p.Validate(); err != nil {
		return err
	}

	got := []byte(p.Password)
	_ = FindOne(&table, &p.Id, &p)
	want := []byte(p.Password)

	return bcrypt.CompareHashAndPassword(want, got)
}

func main() {
	lambda.Start(Handle)
}
