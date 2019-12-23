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

func NewProduct(s string) (Product, error) {
	var p Product
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return p, err
	} else if len(p.Name) < 3 {
		return p, fmt.Errorf("bad name [%s], must be at least 3 characters in length", p.Name)
	} else if p.Price < 0 {
		return p, fmt.Errorf("bad price (integer) [%d]", p.Price)
	} else if p.Id != "" {
		return p, nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return p, err
	} else {
		p.Id = id.String()
		return p, nil
	}
}

var t = os.Getenv("PRODUCT_TABLE")

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	session := r.QueryStringParameters["session"]
	fmt.Printf("REQUEST [%s]: ip=[%s], session=[%s], cmd=[%s], body=[%s]\n", cmd, ip, session, cmd, body)

	switch cmd {

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
		av := r.QueryStringParameters["brand-id"]
		if err := FindAllByAttribute(&t, &an, &av, &p); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "find-by-id":
		var p Product
		s := r.QueryStringParameters["id"]
		if err := FindOne(&t, &s, &p); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "find-by-ids":
		var p []Product
		ss := strings.Split(r.QueryStringParameters["ids"], ",")
		if err := FindAllById(t, ss, &p); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "save":
		if p, err := NewProduct(r.Body); err != nil {
			return BadGateway().Error(err).Build()
		} else if _, err := ValidateSession(session, ip); err != nil {
			return Unauthorized().Error(err).Build()
		} else if err := Put(p, &t); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	default:
		return BadRequest().Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
