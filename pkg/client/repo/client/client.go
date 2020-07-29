package client

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

const deleted = "deleted"

var db *dynamodb.DynamoDB

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}); err != nil {
		log.Fatalf("Failed to connect to AWS: %s", err.Error())
	} else {
		db = dynamodb.New(sess)
	}
}

func FindById(table, id string, i interface{}) error {
	if out, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
	}); err != nil {
		return err
	} else if err := dynamodbattribute.UnmarshalMap(out.Item, &i); err != nil {
		return err
	} else {
		fmt.Println(out)
		return nil
	}
}

func FindByIds(table string, ids []string, i interface{}) error {
	var keys []map[string]*dynamodb.AttributeValue
	for _, s := range ids {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
	}
	if results, err := db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{table: {Keys: keys}}}); err != nil {
		return err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[table], &i); err != nil {
		return err
	} else {
		return nil
	}
}

func FindAll(tableName string, attributes map[string]string, i interface{}) error {
	f := expression.AttributeNotExists(expression.Name(deleted))
	for k, v := range attributes {
		f = f.And(expression.Name(k).Equal(expression.Value(v)))
	}
	if exp, err := expression.NewBuilder().WithFilter(f).Build(); err != nil {
		return err
	} else if out, err := db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		ProjectionExpression:      exp.Projection(),
		TableName:                 aws.String(tableName),
	}); err != nil {
		return err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &i); err != nil {
		return err
	} else {
		return nil
	}
}

func Save(tableName string, i interface{}) error {

	item, err := dynamodbattribute.MarshalMap(&i)
	if err != nil {
		return err
	}

	nowStr := time.Now().UTC().Format(time.RFC3339)
	nowVal := &dynamodb.AttributeValue{S: &nowStr}
	item["modified"] = nowVal

	if id, ok := item["id"]; !ok || id.S == nil {
		s, _ := uuid.NewUUID()
		item["id"] = &dynamodb.AttributeValue{S: aws.String(s.String())}
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{Item: item, TableName: aws.String(tableName)})
	return err
}

func Remove(tableName, id string, i interface{}) error {
	if err := FindById(tableName, id, i); err != nil {
		return nil
	}
	item, err := dynamodbattribute.MarshalMap(&i)
	if err != nil {
		return err
	}
	nowStr := time.Now().UTC().Format(time.RFC3339)
	nowVal := &dynamodb.AttributeValue{S: &nowStr}
	item["deleted"] = nowVal
	_, err = db.PutItem(&dynamodb.PutItemInput{Item: item, TableName: aws.String(tableName)})
	return err
}

func Update(tableName, id, keyword string, ids []string) error {
	_, err := db.UpdateItem(&dynamodb.UpdateItemInput{
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 aws.String(tableName),
		Key:                       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
		UpdateExpression:          aws.String(keyword + " :p"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":p": {SS: aws.StringSlice(ids)}},
	})
	return err
}
