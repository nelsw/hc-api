package repo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"hc-api/model"
	"hc-api/repo/dynamo"
	"os"
	"strings"
)

var userTable = os.Getenv("USER_TABLE")

// Finds user by email address (PK).
func FindUserByEmail(s *string) (user *model.User, err error) {
	if result, err := dynamo.GetItem(userKey(s), &userTable); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

// Updates the specified attributes of a user entity.
func UpdateUser(k, v, e *string) error {
	return dynamo.Update(&dynamodb.UpdateItemInput{
		Key:              userKey(k),
		UpdateExpression: e,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				SS: aws.StringSlice([]string{*v}),
			},
		},
	})
}

// Returns the simple key for retrieving a user entity
func userKey(s *string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"email": {
			S: aws.String(strings.ToLower(*s)),
		},
	}
}
