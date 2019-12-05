package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"hc-api/model"
	"hc-api/service/dynamo"
	"os"
)

var userTable = os.Getenv("USER_TABLE")
var userPasswordTable = os.Getenv("USER_PASSWORD_TABLE")
var userEmailTable = os.Getenv("USER_EMAIL_TABLE")

// Returns the simple key for retrieving a user entity
func key(s *string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"id": {
			S: s,
		},
	}
}

// Finds user by id (PK).
func FindUserById(s *string) (user *model.User, err error) {
	if result, err := dynamo.GetItem(key(s), &userTable); err == nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	}
	return user, nil
}

// Returns the user password entity by providing the user password primary key.
func FindUserPasswordById(s *string) (up *model.UserPassword, err error) {
	if result, err := dynamo.GetItem(key(s), &userPasswordTable); err == nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &up)
	}
	return up, err
}

// Returns the user email entity by providing the user email primary key.
func FindUserEmailById(s *string) (ue *model.UserEmail, err error) {
	if result, err := dynamo.GetItem(key(s), &userEmailTable); err == nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &ue)
	}
	return ue, err
}

// Updates the specified attributes of a user entity.
func UpdateUser(k, e *string, v *[]string) error {
	return dynamo.Update(&dynamodb.UpdateItemInput{
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        &userTable,
		Key:              key(k),
		UpdateExpression: e,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				SS: aws.StringSlice(*v),
			},
		},
	})
}
