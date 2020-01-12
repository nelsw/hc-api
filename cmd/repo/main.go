package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"hc-api/pkg/apigw"
	"log"
	"os"
	"time"
)

type RepoRequest struct {
	Command    string        `json:"command"`
	Body       string        `json:"body"`
	Data       interface{}   `json:"data,omitempty"`
	DataArr    []interface{} `json:"data_arr,omitempty"`
	Table      string        `json:"table"`
	Expression string        `json:"expression,omitempty"`
	Id         string        `json:"id,omitempty"`
	Ids        []string      `json:"ids,omitempty"`
	BrandId    string        `json:"brand_id,omitempty"`
}

const deleted = "deleted"
const brandId = "brand_id"

var db *dynamodb.DynamoDB

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	}); err != nil {
		log.Fatalf("Failed to connect to AWS: %s", err.Error())
	} else {
		db = dynamodb.New(sess)
	}
}

func HandleRequest(r RepoRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST   [%v]\n", r)

	switch r.Command {

	case "find-all":
		f := expression.AttributeNotExists(expression.Name(deleted))
		if exp, err := expression.NewBuilder().WithFilter(f).Build(); err != nil {
			return apigw.BadRequest(err)
		} else if out, err := db.Scan(&dynamodb.ScanInput{
			ExpressionAttributeNames:  exp.Names(),
			ExpressionAttributeValues: exp.Values(),
			FilterExpression:          exp.Filter(),
			ProjectionExpression:      exp.Projection(),
			TableName:                 &r.Table,
		}); err != nil {
			return apigw.BadRequest(err)
		} else if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &r.DataArr); err != nil {
			return apigw.BadRequest(err)
		} else {
			return apigw.Ok(&r.DataArr)
		}

	case "find-by-brand-id":
		f := expression.Name(brandId).Equal(expression.Value(&r.BrandId))
		if exp, err := expression.NewBuilder().WithFilter(f).Build(); err != nil {
			return apigw.BadRequest(err)
		} else if out, err := db.Scan(&dynamodb.ScanInput{
			ExpressionAttributeNames:  exp.Names(),
			ExpressionAttributeValues: exp.Values(),
			FilterExpression:          exp.Filter(),
			ProjectionExpression:      exp.Projection(),
			TableName:                 &r.Table,
		}); err != nil {
			return apigw.BadRequest(err)
		} else if err := dynamodbattribute.UnmarshalListOfMaps(out.Items, &r.DataArr); err != nil {
			return apigw.BadRequest(err)
		} else {
			return apigw.Ok(&r.DataArr)
		}

	case "find-by-id":
		if out, err := db.GetItem(&dynamodb.GetItemInput{
			TableName: &r.Table,
			Key:       map[string]*dynamodb.AttributeValue{"id": {S: &r.Id}},
		}); err != nil {
			return apigw.BadRequest(err)
		} else if err := dynamodbattribute.UnmarshalMap(out.Item, &r.Data); err != nil {
			return apigw.BadRequest(err)
		} else {
			return apigw.Ok(&r.Data)
		}

	case "find-by-ids":
		var keys []map[string]*dynamodb.AttributeValue
		for _, s := range r.Ids {
			keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
		}
		if results, err := db.BatchGetItem(&dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{r.Table: {Keys: keys}}}); err != nil {
			return apigw.BadRequest(err)
		} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[r.Table], &r.DataArr); err != nil {
			return apigw.BadRequest(err)
		} else {
			return apigw.Ok(&r.DataArr)
		}

	case "save":
		if err := json.Unmarshal([]byte(r.Body), &r.Data); err != nil {
			return apigw.BadRequest(err)
		} else if item, err := dynamodbattribute.MarshalMap(&r.Data); err != nil {
			return apigw.BadRequest(err)
		} else {
			if id, ok := item["id"]; !ok || id.S == nil {
				s, _ := uuid.NewUUID()
				item["id"] = &dynamodb.AttributeValue{S: aws.String(s.String())}
			}
			item["modified"] = &dynamodb.AttributeValue{S: aws.String(time.Now().UTC().Format(time.RFC3339))}
			if _, err := db.PutItem(&dynamodb.PutItemInput{Item: item, TableName: &r.Table}); err != nil {
				return apigw.BadRequest(err)
			} else {
				return apigw.Ok(&r.Data)
			}
		}

	case "update":
		if _, err := db.UpdateItem(&dynamodb.UpdateItemInput{
			ReturnValues:              aws.String("UPDATED_NEW"),
			TableName:                 &r.Table,
			Key:                       map[string]*dynamodb.AttributeValue{"id": {S: &r.Id}},
			UpdateExpression:          &r.Expression,
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":p": {SS: aws.StringSlice(r.Ids)}},
		}); err != nil {
			return apigw.BadRequest(err)
		} else {
			return apigw.Ok()
		}
	}

	return apigw.BadRequest()
}

func main() {
	lambda.Start(HandleRequest)
}
