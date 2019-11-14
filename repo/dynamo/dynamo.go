// This package is responsible for exporting generic dynamo methods to domain specific repository ƒ's.
package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

var db *dynamodb.DynamoDB

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}); err != nil {
		log.Fatalf("Failed to connect to AWS: %s", err.Error())
	} else {
		db = dynamodb.New(sess)
	}
}

func Scan(tableName *string) (*dynamodb.ScanOutput, error) {
	return db.Scan(&dynamodb.ScanInput{TableName: tableName})
}

func GetItem(key map[string]*dynamodb.AttributeValue, tableName *string) (*dynamodb.GetItemOutput, error) {
	if output, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: tableName,
		Key:       key,
	}); err != nil {
		return nil, err
	} else {
		return output, err
	}
}

func GetBatch(keys []map[string]*dynamodb.AttributeValue, tableName string) (*dynamodb.BatchGetItemOutput, error) {
	if output, err := db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {Keys: keys},
		},
	}); err != nil {
		return nil, err
	} else {
		return output, nil
	}
}

func PutItem(item map[string]*dynamodb.AttributeValue, tableName *string) error {
	if _, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func Put(v interface{}, s *string) error {
	if item, err := dynamodbattribute.MarshalMap(&v); err != nil {
		return err
	} else if _, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: s,
	}); err != nil {
		return err
	} else {
		return nil
	}
}
