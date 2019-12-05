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
	response "github.com/nelsw/hc-util/aws"
	"hc-api/service"
	"net/http"
	"os"
	"strings"
)

var productTable = os.Getenv("PRODUCT_TABLE")

type Product struct {
	Id            string   `json:"id"`
	Session       string   `json:"session"`
	Sku           string   `json:"sku"`
	Category      string   `json:"category"`
	Subcategory   string   `json:"subcategory"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PriceInteger  string   `json:"price_integer"`
	PriceFraction string   `json:"price_fraction"`
	Quantity      string   `json:"quantity"`
	Unit          string   `json:"unit"`
	Owner         string   `json:"owner"`
	ImageSet      []string `json:"image_set"`
	ShipFrom      string   `json:"ship_from"`
}

func (p *Product) Unmarshal(s string) error {
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return err
	} else if len(p.Name) < 3 {
		return fmt.Errorf("bad name [%s], must be at least 3 characters in length", p.Name)
	} else if p.PriceInteger == "" {
		return fmt.Errorf("bad price (integer) [%s]", p.PriceInteger)
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
	if result, err := service.Scan(&productTable); err != nil {
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

func findAllProductsByIds(ss *[]string) (*[]Product, error) {
	var pp []Product
	var keys []map[string]*dynamodb.AttributeValue
	for _, s := range *ss {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
	}
	if results, err := service.GetBatch(keys, productTable); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[productTable], &pp); err != nil {
		return nil, err
	} else {
		return &pp, nil
	}
}

func saveProduct(p *Product) error {
	p.Session = ""
	return service.Put(p, &productTable)
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	fmt.Printf("REQUEST [%s]: [%v]", cmd, r)

	switch cmd {

	case "save":
		var p Product
		if err := p.Unmarshal(r.Body); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if _, err := service.ValidateSession(p.Session, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := saveProduct(&p); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&p).Build()
		}

	case "find-all":
		if products, err := findAllProducts(); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&products).Build()
		}

	case "find-by-ids":
		csv := r.QueryStringParameters["ids"]
		ids := strings.Split(csv, ",")
		if products, err := findAllProductsByIds(&ids); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&products).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
