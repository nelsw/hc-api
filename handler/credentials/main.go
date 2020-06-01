package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	lambdo "github.com/aws/aws-sdk-go/service/lambda"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"sam-app/pkg/factory/apigwp"
)

var (
	l     *lambdo.Lambda
	db    *dynamodb.DynamoDB
	table = os.Getenv("CREDENTIALS_TABLE")
)

type Credentials struct {
	Username string `json:"username"` // email for v1
	Password string `json:"password"`
}

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}); err != nil {
		log.Fatalf("Failed to connect to AWS: %s", err.Error())
	} else {
		db = dynamodb.New(sess)
	}
}

func (e *Credentials) Validate() error {
	p := e.Password
	if e.Username == "" {
		return fmt.Errorf("bad username [%s]", e.Username)
	} else if e.Password == "" {
		return fmt.Errorf("bad password [%s]", e.Password)
	} else if out, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: &table,
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(e.Username)}},
	}); err != nil {
		return err
	} else if err := dynamodbattribute.UnmarshalMap(out.Item, &e); err != nil {
		return err
	} else if err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(p)); err != nil {
		return err
	} else {
		return nil
	}
}

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var e Credentials
	if _, err := apigwp.Request(r, &e); err != nil {
		return apigwp.Response(400, err)
	}
	// token stuff
	return apigwp.Response(200, "no case")
}

func main() {
	lambda.Start(Handle)
}
