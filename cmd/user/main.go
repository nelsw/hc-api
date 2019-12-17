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

var table = os.Getenv("USER_TABLE")

// Primary user object for the domain, visible to client and server.
// Each property is a reference to a unique ID, or collection of unique ID's.
type User struct {
	Id         string   `json:"id,omitempty"`
	ProfileId  string   `json:"profile_id,omitempty"`
	AddressIds []string `json:"address_ids,omitempty"`
	ProductIds []string `json:"product_ids,omitempty"`
	OrderIds   []string `json:"order_ids,omitempty"`
	SaleIds    []string `json:"sale_ids,omitempty"`
	Session    string   `json:"session"`
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]", cmd, ip, body)

	switch cmd {

	case "login":
		var u User
		if id, err := service.VerifyCredentials(body); err != nil {
			return response.New().Code(http.StatusBadGateway).Text(err.Error()).Build()
		} else if result, err := service.Get(&table, &id); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if err := dynamodbattribute.UnmarshalMap(result.Item, &u); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if cookie, err := service.NewSession(u.Id, ip); err != nil {
			return response.New().Code(http.StatusInternalServerError).Build()
		} else {
			u.Session = cookie
			u.Id = "" // keep user id hidden
			return response.New().Code(http.StatusOK).Toke(cookie).Data(&u).Build()
		}

	case "find":
		var u User
		if id, err := service.ValidateSession(u.Session, ip); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := service.FindOne(&table, &id, &u); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&u).Build()
		}

	case "update":
		var u service.SliceUpdate
		if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if id, err := service.ValidateSession(u.Session, ip); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := service.UpdateSlice(&id, &u.Expression, &table, &u.Val); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Build()
		}

	case "register":
		// todo - create User, UserPassword, and UserProfile entities ... also verify email address.
		return response.New().Code(http.StatusNotImplemented).Build()

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
