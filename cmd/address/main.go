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

var addressTable = os.Getenv("ADDRESS_TABLE")

type Address struct {
	Id      string `json:"id"`
	Session string `json:"session,omitempty"`
	Street1 string `json:"street_1,omitempty"`
	Street2 string `json:"street_2,omitempty"`
	UnitNum string `json:"unit_num,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip5    string `json:"zip_5,omitempty"`
	Zip4    string `json:"zip_4,omitempty"`
}

func (a *Address) Unmarshal(s string) error {
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		return err
	} else if a.Street1 == "" {
		return fmt.Errorf("bad street [%s]", a.Street1)
	} else if a.City == "" {
		return fmt.Errorf("bad city [%s]", a.City)
	} else if a.State == "" {
		return fmt.Errorf("bad state [%s]", a.State)
	} else if a.Zip5 == "" {
		return fmt.Errorf("bad zip [%s]", a.Zip5)
	} else if a.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		a.Id = id.String()
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
	if results, err := service.GetBatch(keys, addressTable); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[addressTable], &aa); err != nil {
		return nil, err
	} else {
		return &aa, nil
	}
}

// Saves an address, creates if new, else updates.
func saveAddress(a *Address) error {
	a.Session = ""
	return service.Put(a, &addressTable)
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	fmt.Printf("REQUEST [%s]: [%v]", cmd, r)

	switch cmd {

	case "save":
		var a Address
		if err := a.Unmarshal(r.Body); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if _, err := service.ValidateSession(a.Session, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := saveAddress(&a); err != nil {
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
