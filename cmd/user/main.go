package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/service"
	. "hc-api/service"
	"net/http"
	"os"
)

var table = os.Getenv("USER_TABLE")

// Primary user object for the domain, visible to client and server.
// Each property is a reference to a unique ID, or collection of unique ID's.
type User struct {
	Id         string   `json:"id"`
	ProfileId  string   `json:"profile_id"`
	AddressIds []string `json:"address_ids"`
	ProductIds []string `json:"product_ids"`
	OrderIds   []string `json:"order_ids"`
	SaleIds    []string `json:"sale_ids"`
	Session    string   `json:"session"`
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	session := r.QueryStringParameters["session"]
	fmt.Printf("REQUEST [%s]: ip=[%s], session=[%s], cmd=[%s], body=[%s]\n", cmd, ip, session, cmd, body)

	switch cmd {

	case "login":
		var u User
		if id, err := service.VerifyCredentials(body); err != nil {
			return BadGateway().Error(err).Build()
		} else if err := service.FindOne(&table, &id, &u); err != nil {
			return BadRequest().Error(err).Build()
		} else if cookie, err := service.NewSession(u.Id, ip); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			u.Session = cookie
			u.Id = "" // keep user id hidden
			return Ok().Data(&u).Build()
		}

	case "find":
		var u User
		session := r.QueryStringParameters["session"]
		if id, err := service.ValidateSession(session, ip); err != nil {
			return Unauthorized().Error(err).Build()
		} else if err := service.FindOne(&table, &id, &u); err != nil {
			return BadRequest().Error(err).Build()
		} else {
			return Ok().Data(&u).Build()
		}

	case "update":
		var u service.SliceUpdate
		if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
			return BadRequest().Error(err).Build()
		} else if id, err := service.ValidateSession(u.Session, ip); err != nil {
			return Unauthorized().Error(err).Build()
		} else if err := service.UpdateSlice(&id, &u.Expression, &table, &u.Val); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Build()
		}

	case "register":
		// todo - create User, UserPassword, and UserProfile entities ... also verify email address.
		return New().Code(http.StatusNotImplemented).Build()

	default:
		return BadRequest().Data(r).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
