package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	. "hc-api/service"
	"os"
	"time"
)

type Offer struct {
	Id           string `json:"id"`
	UserId       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	ProductId    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductPrice int64  `json:"product_price"`
	ProductQty   int    `json:"product_qty"`
	ProductImg   string `json:"product_img,omitempty"`
	Total        int64  `json:"total"`
	Created      string `json:"created"`
}

var tableName = os.Getenv("OFFER_TABLE")

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := request.QueryStringParameters["cmd"]
	body := request.Body
	ip := request.RequestContext.Identity.SourceIP
	session := request.QueryStringParameters["session"]
	fmt.Printf("REQUEST cmd=[%s], ip=[%s], session=[%s], body=[%s]\n", cmd, ip, session, body)

	switch request.QueryStringParameters["cmd"] {

	case "save":

		var offer Offer
		if err := json.Unmarshal([]byte(request.Body), &offer); err != nil {
			return BadGateway().Error(err).Build()
		}

		ip := request.RequestContext.Identity.SourceIP
		session := request.QueryStringParameters["session"]

		um, err := Invoke().Handler("User").IP(ip).CMD("find").Session(session).Build()
		if err != nil {
			return BadRequest().Error(err).Build()
		}

		id := fmt.Sprintf("%v", um["profile_id"])
		pm, err := Invoke().Handler("UserProfile").IP(ip).Session(session).CMD("find").QSP("id", id).Build()
		if err != nil {
			return BadRequest().Error(err).Build()
		}

		offer.Email = fmt.Sprintf("%v", pm["email"])
		offer.Phone = fmt.Sprintf("%v", pm["phone"])
		offer.FirstName = fmt.Sprintf("%v", pm["first_name"])
		offer.LastName = fmt.Sprintf("%v", pm["last_name"])
		offer.Created = time.Now().UTC().Format(time.RFC3339)
		offer.Total = int64(offer.ProductQty) * offer.ProductPrice

		if err := Put(offer, &tableName); err != nil {
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
