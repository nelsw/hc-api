package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	. "hc-api/service"
	"os"
)

var t = os.Getenv("USER_TABLE")

// Primary user object for the domain, visible to client and server.
// Each property is a reference to a unique ID, or collection of unique ID's.
type User struct {
	Id         string   `json:"id"`
	ProfileId  string   `json:"profile_id"`
	AddressIds []string `json:"address_ids"`
	ProductIds []string `json:"product_ids"`
	OrderIds   []string `json:"order_ids,omitempty"`
	OfferIds   []string `json:"sale_ids,omitempty"`
	SaleIds    []string `json:"sale_ids,omitempty"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := request.QueryStringParameters["cmd"]
	body := request.Body
	ip := request.RequestContext.Identity.SourceIP
	session := request.QueryStringParameters["session"]
	fmt.Printf("REQUEST cmd=[%s], ip=[%s], session=[%s], body=[%s]\n", cmd, ip, session, body)

	switch cmd {

	case "login":
		var u User
		if id, err := VerifyCredentials(request.Body); err != nil {
			return BadGateway().Error(err).Build()
		} else if err := FindOne(&t, &id, &u); err != nil {
			return BadRequest().Error(err).Build()
		} else if cookie, err := NewSession(id, request.RequestContext.Identity.SourceIP); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			u.Id = cookie
			return Ok().Data(&u).Build()
		}

	case "find":
		var u User
		session := request.QueryStringParameters["session"]
		if id, err := ValidateSession(session, request.RequestContext.Identity.SourceIP); err != nil {
			return Unauthorized().Error(err).Build()
		} else if err := FindOne(&t, &id, &u); err != nil {
			return NotFound().Error(err).Build()
		} else {
			return Ok().Data(&u).Build()
		}

	case "update":
		var u SliceUpdate
		if err := json.Unmarshal([]byte(request.Body), &u); err != nil {
			return BadRequest().Error(err).Build()
		} else if id, err := ValidateSession(u.Session, request.RequestContext.Identity.SourceIP); err != nil {
			return Unauthorized().Error(err).Build()
		} else if err := UpdateSlice(&id, &u.Expression, &t, &u.Val); err != nil {
			return InternalServerError().Error(err).Build()
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
