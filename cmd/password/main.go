package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
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
		default:
			// do not increment length for unrecognized characters
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

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST [%v]", request)

	switch request.QueryStringParameters["cmd"] {

	case "verify":
		client := request.QueryStringParameters["p"]
		id := request.QueryStringParameters["id"]
		var server Password
		if err := FindOne(&t, &id, &server); err != nil {
			return NotFound().Error(err).Build()
		} else if err := bcrypt.CompareHashAndPassword([]byte(server.Password), []byte(client)); err != nil {
			return Unauthorized().Error(err).Build()
		} else {
			return Ok().Build()
		}

	default:
		return BadRequest().Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
