package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
	// not used, yet.
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	// basic details.
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	ImageSet    []string `json:"image_set"`
	Quantity    string   `json:"quantity"`
	// inches
	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`
	// pounds.ounces
	Weight float32 `json:"weight"`
	// deprecated
	Unit          string `json:"unit"`
	Owner         string `json:"owner"` // user_id
	PriceInteger  string `json:"price_integer"`
	PriceFraction string `json:"price_fraction"`
	ZipCode       string `json:"zip_code"`

	// transient
	Session string `json:"session"`
}

func (p *Product) Unmarshal(s string) error {
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return err
	} else if len(p.Name) < 3 {
		return fmt.Errorf("bad name [%s], must be at least 3 characters in length", p.Name)
	} else if p.Price < 0 {
		return fmt.Errorf("bad price (integer) [%d]", p.Price)
	} else if p.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		p.Id = id.String()
		return nil
	}
}

func findAllProducts() ([]Product, error) {
	if result, err := Scan(&table); err != nil {
		return nil, err
	} else {
		var goods []Product
		for _, item := range result.Items {
			good := Product{}
			if err := dynamodbattribute.UnmarshalMap(item, &good); err != nil {
				return nil, err
			} else {
				goods = append(goods, good)
			}
		}
		return goods, nil
	}
}

func findProductsByIds(ss *[]string) (*[]Product, error) {
	var pp []Product
	var keys []map[string]*dynamodb.AttributeValue
	for _, s := range *ss {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
	}
	if results, err := GetBatch(keys, table); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[table], &pp); err != nil {
		return nil, err
	} else {
		return &pp, nil
	}
}

var table = os.Getenv("PRODUCT_TABLE")

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]", cmd, ip, body)

	switch cmd {

	case "save":
		var p Product
		if err := p.Unmarshal(r.Body); err != nil {
			return BadGateway().Error(err).Build()
		} else if _, err := ValidateSession(p.Session, ip); err != nil {
			return Unauthorized().Error(err).Build()
		} else if err := Put(p, &table); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "find-all":
		if products, err := findAllProducts(); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&products).Build()
		}

	case "find":
		var p Product
		id := r.QueryStringParameters["id"]
		if err := FindOne(&table, &id, &p); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "find-by-ids":
		csv := r.QueryStringParameters["ids"]
		ids := strings.Split(csv, ",")
		if products, err := findProductsByIds(&ids); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&products).Build()
		}

	default:
		return BadRequest().Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
