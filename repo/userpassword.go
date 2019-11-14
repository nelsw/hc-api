package repo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"hc-api/model"
	"hc-api/repo/dynamo"
	"os"
)

var userPasswordTable = os.Getenv("USER_PASSWORD_TABLE")

func FindUserPasswordById(s *string) (up *model.UserPassword, err error) {
	if result, err := dynamo.GetItem(userPasswordKey(s), &userPasswordTable); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalMap(result.Item, up); err != nil {
		return nil, err
	} else {
		return up, nil
	}
}

func SaveUserPassword(up *model.UserPassword) error {
	return dynamo.Put(up, &userPasswordTable)
}

func userPasswordKey(s *string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{"id": {S: s}}
}
