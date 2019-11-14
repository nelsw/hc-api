package repo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"hc-api/model"
	"hc-api/repo/dynamo"
	"os"
)

var addressTable = os.Getenv("ADDRESS_TABLE")

// Finds address by id (PK).
func FindAddress(address *model.Address) error {
	if result, err := dynamo.GetItem(addressKey(address), &addressTable); err != nil {
		return err
	} else {
		return dynamodbattribute.UnmarshalMap(result.Item, &address)
	}
}

// Finds addresses by each address id (PK).
func FindAllAddresses(aa *[]model.Address) error {
	var keys []map[string]*dynamodb.AttributeValue
	for _, a := range *aa {
		keys = append(keys, addressKey(&a))
	}
	if results, err := dynamo.GetBatch(keys, addressTable); err != nil {
		return err
	} else {
		return dynamodbattribute.UnmarshalListOfMaps(results.Responses[addressTable], &aa)
	}
}

// Saves an address, creates if new, else updates.
func SaveAddress(address *model.Address) error {
	if item, err := dynamodbattribute.MarshalMap(&address); err != nil {
		return err
	} else {
		return dynamo.PutItem(item, &addressTable)
	}
}

// Returns the simple key for retrieving an address entity
func addressKey(address *model.Address) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: &address.Id,
		},
	}
}
