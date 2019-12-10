package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/service"
	"net/http"
	"os"
	"strings"
)

var table = os.Getenv("ADDRESS_TABLE")

type Address struct {
	Id      string `json:"id"`
	Session string `json:"session"`
	Street  string `json:"street"`
	Unit    string `json:"unit,omitempty"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip5    string `json:"zip_5"`
	Zip4    string `json:"zip_4,omitempty"`
}

func (a *Address) Validate() error {
	if b, err := json.Marshal(a); err != nil {
		return err
	} else if str, err := service.VerifyAddress(string(b)); err != nil {
		return err
	} else if err := json.Unmarshal([]byte(str), &a); err != nil {
		return err
	} else {
		return nil
	}
}

// Finds addresses by each address id (PK).
func findAllAddressesByIds(ss *[]string) (*[]Address, error) {
	var aa []Address
	var keys []map[string]*dynamodb.AttributeValue
	for _, s := range *ss {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
	}
	if results, err := service.GetBatch(keys, table); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[table], &aa); err != nil {
		return nil, err
	} else {
		return &aa, nil
	}
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]", cmd, ip, body)

	switch cmd {

	case "save":
		var a Address
		if err := json.Unmarshal([]byte(body), &a); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if _, err := service.ValidateSession(a.Session, ip); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if b, err := json.Marshal(a); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if str, err := service.VerifyAddress(string(b)); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := json.Unmarshal([]byte(str), &a); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := service.Put(a, &table); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&a).Build()
		}

	case "find-by-ids":
		csv := r.QueryStringParameters["ids"]
		ids := strings.Split(csv, ",")
		if aa, err := findAllAddressesByIds(&ids); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&aa).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
