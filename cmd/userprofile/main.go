// UserProfile is exactly what it appears to be, a user profile domain model entity.
// It also promotes separation of concerns by decoupling user profile details from the primary user entity. IF
// UserProfile.EmailOld != UserProfile.EmailNew, AND User.Email == UserProfile.EmailOld, THEN we must prompt the user to
// confirm new email address. IF UserProfile.Password1 is not blank AND UserProfile.Password2 is not blank AND valid AND
// UserProfile.Password1 == UserProfile.Password2, then we update the UserPassword entity and return OK.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/service"
	"hc-api/service/dynamo"
	"net/http"
	"os"
)

var table = os.Getenv("USER_PROFILE_TABLE")

type UserProfile struct {
	Id        string `json:"id"`
	Session   string `json:"session"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Password1 string `json:"password_1,omitempty"`
	Password2 string `json:"password_2,omitempty"`
}

func (up *UserProfile) Unmarshal(s string) error {
	if err := json.Unmarshal([]byte(s), &up); err != nil {
		return err
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		up.Id = id.String()
		return nil
	}
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	fmt.Printf("REQUEST [%s]: [%v]", cmd, r)

	switch cmd {

	case "save":
		var up UserProfile
		if err := up.Unmarshal(r.Body); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if _, err := service.ValidateSession(up.Session, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := dynamo.Put(up, &table); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&up).Build()
		}

	case "find":
		var up UserProfile
		id := r.QueryStringParameters["id"]
		if result, err := dynamo.Get(&table, &id); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if err := dynamodbattribute.UnmarshalMap(result.Item, &up); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&up).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
