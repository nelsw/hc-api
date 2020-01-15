package repository

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"hc-api/internal/class"
	"time"
)

const deleted = "deleted"
const brandId = "brand_id"

func FindByIds(r class.Request) ([]byte, error) {
	var keys []map[string]*dynamodb.AttributeValue
	for _, s := range r.Ids {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
	}
	if results, err := db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{r.Table: {Keys: keys}}}); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[r.Table], &r.DataArr); err != nil {
		return nil, err
	} else {
		return json.Marshal(&r.DataArr)
	}
}

func FindById(r class.Request) ([]byte, error) {
	if out, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: &r.Table,
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: &r.Ids[0]}},
	}); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalMap(out.Item, &r.Data); err != nil {
		return nil, err
	} else {
		return json.Marshal(&r.Data)
	}
}

func FindByBrandId(r class.Request) ([]byte, error) {
	f := expression.Name(brandId).Equal(expression.Value(&r.BrandId))
	if exp, err := expression.NewBuilder().WithFilter(f).Build(); err != nil {
		return nil, err
	} else if out, err := db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		ProjectionExpression:      exp.Projection(),
		TableName:                 &r.Table,
	}); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &r.DataArr); err != nil {
		return nil, err
	} else {
		return json.Marshal(&r.DataArr)
	}
}

func FindAll(r class.Request) ([]byte, error) {
	f := expression.AttributeNotExists(expression.Name(deleted))
	if exp, err := expression.NewBuilder().WithFilter(f).Build(); err != nil {
		return nil, err
	} else if out, err := db.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		ProjectionExpression:      exp.Projection(),
		TableName:                 &r.Table,
	}); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &r.DataArr); err != nil {
		return nil, err
	} else {
		return json.Marshal(&r.DataArr)
	}
}

func Save(r class.Request) ([]byte, error) {
	item, err := dynamodbattribute.MarshalMap(&r.Data)
	if err != nil {
		return nil, err
	}

	nowStr := time.Now().UTC().Format(time.RFC3339)
	nowVal := &dynamodb.AttributeValue{S: &nowStr}
	item["modified"] = nowVal
	if id, ok := item["modified"]; !ok || id.S == nil {
		item["created"] = nowVal
	}
	if id, ok := item["id"]; !ok || id.S == nil {
		s, _ := uuid.NewUUID()
		item["id"] = &dynamodb.AttributeValue{S: aws.String(s.String())}
	}

	if _, err := db.PutItem(&dynamodb.PutItemInput{Item: item, TableName: &r.Table}); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalMap(item, &r.Data); err != nil {
		return nil, err
	} else {
		return json.Marshal(&r.Data)
	}
}

func Update(r class.Request) ([]byte, error) {
	if _, err := db.UpdateItem(&dynamodb.UpdateItemInput{
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 &r.Table,
		Key:                       map[string]*dynamodb.AttributeValue{"id": {S: &r.Ids[0]}},
		UpdateExpression:          &r.Expression,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":p": {SS: aws.StringSlice(r.Ids)}},
	}); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}
