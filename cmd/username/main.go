package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/service"
	"net/http"
	"os"
)

var table = os.Getenv("USERNAME_TABLE")

// Used for login and registration use cases.
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]", cmd, ip, body)

	switch cmd {

	case "verify":
		var uc UserCredentials
		var un Username
		if err := json.Unmarshal([]byte(body), &uc); err != nil {
			return response.New().Code(http.StatusBadGateway).Text(err.Error()).Build()
		} else if result, err := service.Get(&table, &uc.Email); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if err := dynamodbattribute.UnmarshalMap(result.Item, &un); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if err := service.VerifyPassword(uc.Password, un.PasswordId); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Text(un.UserId).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
