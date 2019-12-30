package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	. "hc-api/service"
	"os"
	"strings"
)

type Product struct {
	// ids and so forth.
	Id        string `json:"id"`
	Sku       string `json:"sku"`
	AddressId string `json:"address_id"`
	BrandId   string `json:"brand_id"`
	Owner     string `json:"owner"` // user_id
	// basic details.
	Category    string   `json:"category"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	ImageSet    []string `json:"image_set"`
	Quantity    string   `json:"quantity"`
	Stock       string   `json:"stock"`
	// packaging details (calc shipping rates)
	Unit   string  `json:"unit"` // LB
	Weight float32 `json:"weight"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Length int     `json:"length"`
}

var t = os.Getenv("PRODUCT_TABLE")

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST [%v]", request)

	switch request.QueryStringParameters["cmd"] {

	case "find-all":
		var p []Product
		if err := FindAll(&t, &p); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "find-by-brand-id":
		var p []Product
		an := "brand_id"
		av := request.QueryStringParameters["brand-id"]
		if err := FindAllByAttribute(&t, &an, &av, &p); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "find-by-id":
		var p Product
		s := request.QueryStringParameters["id"]
		if err := FindOne(&t, &s, &p); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "find-by-ids":
		var p []Product
		ss := strings.Split(request.QueryStringParameters["ids"], ",")
		if err := FindAllById(t, ss, &p); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "save":
		var p Product
		if err := json.Unmarshal([]byte(request.Body), &p); err != nil {
			return BadGateway().Error(err).Build()
		} else if len(p.Name) < 3 {
			return BadRequest().Error(fmt.Errorf("bad name [%s], must be at least 3 characters in length", p.Name)).Build()
		} else if p.Price < 0 {
			return BadRequest().Error(fmt.Errorf("bad price (integer) [%d]", p.Price)).Build()
		}

		ip := request.RequestContext.Identity.SourceIP
		session := request.QueryStringParameters["session"]
		if userId, err := ValidateSession(session, ip); err != nil {
			return Unauthorized().Error(err).Build()
		} else {
			p.Owner = userId
			if p.Id == "" {
				id, _ := uuid.NewUUID()
				p.Id = id.String()
			}
			if err := Put(p, &t); err != nil {
				return InternalServerError().Error(err).Build()
			} else {
				return Ok().Data(&p).Build()
			}
		}

	default:
		return BadRequest().Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
