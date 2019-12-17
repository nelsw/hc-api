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
// val, slice of ids to set
// session, provides info for transaction
// expression, eg. "set package_ids = :p" other expression commands:
// SET - modify or add item attributes
// REMOVE - delete attributes from an item
// ADD - update numbers and sets
// DELETE - remove elements from a set
type SliceUpdate struct {
	Val        []string `json:"val"`
	Expression string   `json:"expression"`
	Session    string   `json:"session"` // valid session required
}

var db *dynamodb.DynamoDB

// todo - cross region initialization via env var
func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}); err != nil {
		log.Fatalf("Failed to connect to AWS: %s", err.Error())
	} else {
		db = dynamodb.New(sess)
	}
}

// all pk column names are identical in our business domain data model
func key(id *string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{"id": {S: id}}
}

// until pagination becomes a requirement, scans are not hurting performance
func Scan(s *string) (*dynamodb.ScanOutput, error) {
	return db.Scan(&dynamodb.ScanInput{TableName: s})
}

// similar to an ORM, this method returns a single entity by providing a table name and pk
func Get(tn, id *string) (*dynamodb.GetItemOutput, error) {
	return db.GetItem(&dynamodb.GetItemInput{TableName: tn, Key: key(id)})
}

// similar to an ORM, this method returns a single entity by providing a table name and pk
func FindOne(tn, id *string, v interface{}) error {
	out, err := db.GetItem(&dynamodb.GetItemInput{TableName: tn, Key: key(id)})
	if err != nil {
		return err
	} else if err := dynamodbattribute.UnmarshalMap(out.Item, &v); err != nil {
		return err
	} else {
		return nil
	}
}

// similar to Get(tn, id), this method returns a batch of entities by providing a table name and pks
func GetBatch(keys []map[string]*dynamodb.AttributeValue, tableName string) (*dynamodb.BatchGetItemOutput, error) {
	return db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{tableName: {Keys: keys}}})
}

// similar to merge or save, this method will only insert and update missing values
func Put(v interface{}, s *string) error {
	if item, err := dynamodbattribute.MarshalMap(&v); err != nil {
		return err
	} else {
		delete(item, "session")
		_, err := db.PutItem(&dynamodb.PutItemInput{Item: item, TableName: s})
		return err
	}
}

// like put but with a dynamodb condition expression
func PutConditionally(v interface{}, s, c *string, e map[string]*dynamodb.AttributeValue) error {
	if item, err := dynamodbattribute.MarshalMap(&v); err == nil {
		return err
	} else {
		delete(item, "session")
		in := &dynamodb.PutItemInput{Item: item, TableName: s, ConditionExpression: c, ExpressionAttributeValues: e}
		_, err := db.PutItem(in)
		return err
	}
}

// deletes an entity by providing a table name and pk
func Delete(id, table *string) error {
	_, err := db.DeleteItem(&dynamodb.DeleteItemInput{Key: key(id), TableName: table})
	return err
}

// as our data model leverages dynamodb's string slice, this method provides a means for updating an entities slice
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
