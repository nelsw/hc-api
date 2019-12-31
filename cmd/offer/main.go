package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	. "hc-api/service"
	"os"
	"time"
)

type Offer struct {
	Id               string `json:"id"`
	UserId           string `json:"user_id"`
	ProfileId        string `json:"profile_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	ProductId        string `json:"product_id"`
	ProductImg       string `json:"product_img"`
	ProductName      string `json:"product_name"`
	ProductUnit      string `json:"product_unit"`
	ProductPrice     int64  `json:"product_price"`
	ProductAddressId string `json:"product_address_id"`
	Details          string `json:"details"`
	Created          string `json:"created"`
}

type OfferRequest struct {
	Command string `json:"command"`
	Session string `json:"session"`
	Offer   Offer  `json:"offer"`
}

const f = "REQUEST offer\n\t     IP=[%s]\n\tcommand=[%s]\n\tsession=[%s]\n\t   body=[%v]\n"

var t = os.Getenv("OFFER_TABLE")

func HandleRequest(proxyRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var r OfferRequest
	if err := json.Unmarshal([]byte(proxyRequest.Body), &r); err != nil {
		panic(err)
	}

	ip := proxyRequest.RequestContext.Identity.SourceIP
	fmt.Printf(f, ip, r.Command, r.Session, r.Offer)

	switch r.Command {

	case "save":

		userId, err := Invoke().Handler("Session").Session(r.Session).IP(ip).CMD("validate").Post()
		if err != nil {
			return Unauthorized().Error(err).Build()
		}

		err = Invoke().Handler("User").IP(ip).CMD("find").Session(r.Session).Marshal(&r.Offer)
		if err != nil {
			return BadRequest().Error(err).Build()
		}

		err = Invoke().Handler("Profile").IP(ip).Session(r.Session).CMD("find").ID(r.Offer.ProfileId).Marshal(&r.Offer)
		if err != nil {
			return BadRequest().Error(err).Build()
		}

		id, _ := uuid.NewUUID()
		r.Offer.Id = id.String()
		r.Offer.UserId = userId
		r.Offer.Created = time.Now().UTC().Format(time.RFC3339)

		if err := Put(r.Offer, &t); err != nil {
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
