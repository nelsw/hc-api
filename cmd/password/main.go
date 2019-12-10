package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	response "github.com/nelsw/hc-util/aws"
	"golang.org/x/crypto/bcrypt"
	"hc-api/service"
	"net/http"
	"os"
	"unicode"
)

var table = os.Getenv("USER_PASSWORD_TABLE")

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

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]", cmd, ip, body)

	switch cmd {

	case "verify":
		p := r.QueryStringParameters["p"]
		id := r.QueryStringParameters["id"]
		var server Password
		if result, err := service.Get(&table, &id); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if err = dynamodbattribute.UnmarshalMap(result.Item, &server); err != nil {
			return response.New().Code(http.StatusUnauthorized).Build()
		} else if err := bcrypt.CompareHashAndPassword([]byte(server.Password), []byte(p)); err != nil {
			return response.New().Code(http.StatusUnauthorized).Build()
		} else {
			return response.New().Code(http.StatusOK).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
