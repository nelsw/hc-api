package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
	"hc-api/pkg/apigw"
	. "hc-api/service"
	"os"
	"unicode"
)

var t = os.Getenv("USER_PASSWORD_TABLE")

// Data structure for persisting and retrieving a users (encrypted) password. The user entity maintains a 1-1 FK
// relationship with the UserPassword entity for referential integrity, where user.PasswordId == userPassword.Id.
type Password struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type PasswordRequest struct {
	Session  string   `json:"session"`
	Password Password `json:"password"`
}

// Validates the UserPassword entity by confirming that both the password and id values are valid.
// Most of the following method body is an adaptation of https://stackoverflow.com/a/25840157.
func (up *Password) Validate() error {
	var number, upper, special bool
	length := 0
	for _, c := range up.Password {
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

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	body := request.Body
	fmt.Printf("REQUEST   [%s]\n", body)

	var r PasswordRequest
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return apigw.BadRequest(err)
	} else if err := r.Password.Validate(); err != nil {
		return apigw.BadRequest(err)
	}

	ip := request.RequestContext.Identity.SourceIP
	if _, err := Invoke().Handler("Session").Session(r.Session).IP(ip).CMD("validate").Post(); err != nil {
		return apigw.BadAuth(err)
	}

	var p Password
	if err := FindOne(&t, &r.Password.Id, &p); err != nil {
		return apigw.BadRequest(err)
	} else if err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(r.Password.Password)); err != nil {
		return apigw.BadAuth(err)
	} else {
		return apigw.Ok()
	}
}

func main() {
	lambda.Start(Handle)
}
