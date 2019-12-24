package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	. "hc-api/service"
	"net/http"
	"os"
	"strings"
)

var table = os.Getenv("USERNAME_TABLE")

// Used for login and registration use cases.
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserCredentials(s string) (UserCredentials, error) {
	var p UserCredentials
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return p, err
	} else {
		p.Email = strings.ToLower(p.Email)
		return p, nil
	}
}

// Data structure for securely associating user entities with their username.
type Username struct {
	Id         string `json:"id"` // email address OR username
	UserId     string `json:"user_id"`
	PasswordId string `json:"password_id"`
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	session := r.QueryStringParameters["session"]
	fmt.Printf("REQUEST cmd=[%s], ip=[%s], session=[%s], body=[%s]\n", cmd, ip, session, body)

	switch cmd {

	case "verify":
		var un Username
		if uc, err := NewUserCredentials(body); err != nil {
			return BadGateway().Error(err).Build()
		} else if err := FindOne(&table, &uc.Email, &un); err != nil {
			return New().Code(http.StatusNotFound).Error(err).Build()
		} else if err := VerifyPassword(uc.Password, un.PasswordId); err != nil {
			return Unauthorized().Error(err).Build()
		} else {
			return Ok().Str(un.UserId).Build()
		}

	default:
		return BadRequest().Data(r).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
