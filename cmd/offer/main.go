package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"hc-api/pkg/apigw"
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

var t = os.Getenv("OFFER_TABLE")

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var r OfferRequest
	body := request.Body
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return apigw.BadRequest(err)
	}

	fmt.Printf("REQUEST   [%s]\n", body)

	if r.Command != "save" {
		return apigw.BadRequest("bad command")
	}

	tkn := r.Session
	ip := request.RequestContext.Identity.SourceIP

	if userId, err := Invoke().Handler("Session").IP(ip).Session(tkn).CMD("validate").Post(); err != nil {
		return apigw.BadAuth(err)
	} else if err := Invoke().Handler("User").IP(ip).Session(tkn).CMD("find").Marshal(&r.Offer); err != nil {
		return apigw.BadRequest(err)
	} else if pid := r.Offer.ProfileId; pid == "" {
		return apigw.BadRequest("bad user id")
	} else if err := Invoke().Handler("Profile").IP(ip).Session(tkn).CMD("find").ID(pid).Marshal(&r.Offer); err != nil {
		return apigw.BadRequest(err)
	} else {

		id, _ := uuid.NewUUID()
		r.Offer.Id = id.String()
		r.Offer.UserId = userId
		r.Offer.Created = time.Now().UTC().Format(time.RFC3339)

		if err := Put(r.Offer, &t); err != nil {
			return apigw.BadRequest(err)
		} else {
			return apigw.Ok()
		}
	}
}

func main() {
	lambda.Start(Handle)
}
