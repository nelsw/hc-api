package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/service"
	"os"
	"strings"
	"unicode"
)

// Used for login and registration use cases.
type ClientCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Data structure for securely associating user entities with their credential and password.
type ServerCredential struct {
	Id         string `json:"id"` // email address
	UserId     string `json:"user_id"`
	PasswordId string `json:"password_id"`
}

var (
	table     = os.Getenv("USERNAME_TABLE")
	LengthErr = fmt.Errorf("bad password, must contain 8-24 characters")
	NumberErr = fmt.Errorf("bad password, must contain at least 1 number")
	CaseErr   = fmt.Errorf("bad password, must contain at least 1 uppercase letter")
	CharErr   = fmt.Errorf("bad password, must contain at least 1 special character")
)

// Validates the UserPassword entity by confirming that both the password and id values are valid.
// Most of the following method body is an adaptation of https://stackoverflow.com/a/25840157.
func (cc *ClientCredential) Validate() error {
	cc.Email = strings.ToLower(cc.Email)
	var number, upper, special bool
	length := 0
	for _, c := range cc.Password {
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

func Handle(cc ClientCredential) (string, error) {
	var sc ServerCredential
	if err := cc.Validate(); err != nil {
		return "", err
	} else if err := service.FindOne(&table, &cc.Email, &sc); err != nil {
		return "", err
	} else if err := service.VerifyPassword(sc.PasswordId, cc.Password); err != nil {
		return "", err
	} else {
		return sc.UserId, nil
	}
}

func main() {
	lambda.Start(Handle)
}
