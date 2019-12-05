package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"hc-api/model"
	"hc-api/service/dynamo"
	"os"
)

var addressTable = os.Getenv("ADDRESS_TABLE")

// Finds addresses by each address id (PK).
func FindAllAddressesByIds(ss *[]string) (*[]model.Address, error) {
	var aa []model.Address
	var keys []map[string]*dynamodb.AttributeValue
	for _, s := range *ss {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
	}
	if results, err := dynamo.GetBatch(keys, addressTable); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[addressTable], &aa); err != nil {
		return nil, err
	} else {
		return &aa, nil
	}
}

// Saves an address, creates if new, else updates.
func SaveAddress(a *model.Address) error {
	a.Session = ""
	if item, err := dynamodbattribute.MarshalMap(&a); err != nil {
		return err
	} else {
		return dynamo.PutItem(item, &addressTable)
	}
}
