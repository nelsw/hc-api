// This package is responsible for exporting generic dynamo methods to domain specific repository ƒ's.
package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

// Used to update an existing user item in an Amazon DynamoDB table.
// SET - modify or add item attributes
// REMOVE - delete attributes from an item
// ADD - update numbers and sets
// DELETE - remove elements from a set
type SliceUpdate struct {
	Id         string   `json:"id,omitempty"`
	Val        []string `json:"val"`
	Expression string   `json:"expression"`
	Session    string   `json:"session"`
}

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

func key(id *string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{"id": {S: id}}
}

func Scan(s *string) (*dynamodb.ScanOutput, error) {
	return db.Scan(&dynamodb.ScanInput{TableName: s})
}

func Get(tn, id *string) (*dynamodb.GetItemOutput, error) {
	return db.GetItem(&dynamodb.GetItemInput{TableName: tn, Key: key(id)})
}

func GetBatch(keys []map[string]*dynamodb.AttributeValue, tableName string) (*dynamodb.BatchGetItemOutput, error) {
	return db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{tableName: {Keys: keys}}})
}

func Put(v interface{}, s *string) error {
	if item, err := dynamodbattribute.MarshalMap(&v); err != nil {
		return err
	} else {
		_, err := db.PutItem(&dynamodb.PutItemInput{Item: item, TableName: s})
		return err
	}
}

func Delete(id, table *string) error {
	_, err := db.DeleteItem(&dynamodb.DeleteItemInput{Key: key(id), TableName: table})
	return err
}

func Update(input *dynamodb.UpdateItemInput) error {
	_, err := db.UpdateItem(input)
	return err
}

func UpdateSlice(k, e, t *string, p *[]string) error {
	_, err := db.UpdateItem(&dynamodb.UpdateItemInput{
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 t,
		Key:                       key(k),
		UpdateExpression:          e,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":p": {SS: aws.StringSlice(*p)}},
	})
	return err
}
