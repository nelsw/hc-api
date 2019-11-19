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

func ScanTable(s *string) (*dynamodb.ScanOutput, error) {
	return db.Scan(&dynamodb.ScanInput{TableName: s})
}

func Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return db.Scan(input)
}

func GetItem(key map[string]*dynamodb.AttributeValue, tableName *string) (*dynamodb.GetItemOutput, error) {
	return db.GetItem(&dynamodb.GetItemInput{TableName: tableName, Key: key})
}

func GetBatch(keys []map[string]*dynamodb.AttributeValue, tableName string) (*dynamodb.BatchGetItemOutput, error) {
	return db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {Keys: keys},
		},
	})
}

func PutItem(item map[string]*dynamodb.AttributeValue, tableName *string) error {
	_, err := db.PutItem(&dynamodb.PutItemInput{Item: item, TableName: tableName})
	return err
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

func Update(input *dynamodb.UpdateItemInput) error {
	_, err := db.UpdateItem(input)
	return err
}
