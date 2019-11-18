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
		return user, err
	} else if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return user, err
	} else {
		return user, err
	}
}

// Saves a user, creates if new, else updates.
func SaveUser(user *model.User) error {
	if item, err := dynamodbattribute.MarshalMap(&user); err != nil {
		return err
	} else {
		return dynamo.PutItem(item, &userTable)
	}
}

// todo - add/remove address id
// todo - add/remove product id
// todo - add/remove order id
// todo - add/remove sale id

// Returns the simple key for retrieving a user entity
func userKey(s *string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"email": {
			S: aws.String(strings.ToLower(*s)),
		},
	}
}
