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
var userPasswordTable = os.Getenv("USER_PASSWORD_TABLE")

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
func UpdateUser(k, e *string, v *[]string) error {
	return dynamo.Update(&dynamodb.UpdateItemInput{
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        &userTable,
		Key:              userKey(k),
		UpdateExpression: e,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				SS: aws.StringSlice(*v),
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

// Returns the user password entity by providing the user password key.
func FindUserPasswordById(s *string) (up *model.UserPassword, err error) {
	if result, err := dynamo.GetItem(map[string]*dynamodb.AttributeValue{"id": {S: s}}, &userPasswordTable); err != nil {
		return up, err
	} else if err := dynamodbattribute.UnmarshalMap(result.Item, &up); err != nil {
		return up, err
	} else {
		return up, err
	}
}
