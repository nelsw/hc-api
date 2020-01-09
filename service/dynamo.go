// This package is responsible for exporting generic dynamo methods to domain specific repository ƒ's.
package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"log"
)

// Used to update an existing user item in an Amazon DynamoDB table.
// val, slice of values, typically ids
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

const deleted = "deleted"

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

// similar to merge or save, this method will only insert and update missing values
func Put(v interface{}, s *string) error {
	if item, err := dynamodbattribute.MarshalMap(&v); err != nil {
		return err
	} else {
		delete(item, "session") // deprecated
		if av, ok := item["id"]; !ok || av.S == nil {
			id, _ := uuid.NewUUID()
			item["id"] = &dynamodb.AttributeValue{S: aws.String(id.String())}
		}
		_, err := db.PutItem(&dynamodb.PutItemInput{Item: item, TableName: s})
		return err
	}
}

// similar to an ORM, this method returns a single entity by providing a table name and pk
func FindOne(tn, id *string, v interface{}) error {
	out, err := db.GetItem(&dynamodb.GetItemInput{TableName: tn, Key: map[string]*dynamodb.AttributeValue{"id": {S: id}}})
	if err == nil {
		err = dynamodbattribute.UnmarshalMap(out.Item, &v)
	}
	return err
}

func FindAll(s *string, v interface{}) error {
	f := expression.AttributeNotExists(expression.Name(deleted))
	if exp, err := expression.NewBuilder().WithFilter(f).Build(); err != nil {
		return err
	} else if out, err := db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		ProjectionExpression:      exp.Projection(),
		TableName:                 s,
	}); err != nil {
		return err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &v); err != nil {
		return err
	} else {
		return nil
	}
}

func FindAllById(tn string, ss []string, v interface{}) error {
	var keys []map[string]*dynamodb.AttributeValue
	for _, s := range ss {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
	}
	if results, err := db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{tn: {Keys: keys}}}); err != nil {
		return err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[tn], &v); err != nil {
		return err
	} else {
		return nil
	}
}

// similar to an ORM, this method returns multiple entities by providing a table name and attribute to filter by
func FindAllByAttribute(tn, an, av *string, v interface{}) error {
	f := expression.Name(*an).Equal(expression.Value(av))
	if exp, err := expression.NewBuilder().WithFilter(f).Build(); err != nil {
		return err
	} else if out, err := db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		ProjectionExpression:      exp.Projection(),
		TableName:                 tn,
	}); err != nil {
		return err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &v); err != nil {
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

// as our data model leverages dynamodb's string slice, this method provides a means for updating an entities slice
func UpdateSlice(k, e, t *string, p *[]string) error {
	_, err := db.UpdateItem(&dynamodb.UpdateItemInput{
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 t,
		Key:                       map[string]*dynamodb.AttributeValue{"id": {S: k}},
		UpdateExpression:          e,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":p": {SS: aws.StringSlice(*p)}},
	})
	return err
}
