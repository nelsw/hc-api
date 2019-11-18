package repo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"hc-api/model"
	"hc-api/repo/dynamo"
	"os"
)

var productTable = os.Getenv("PRODUCT_TABLE")

func FindAllProducts() ([]model.Product, error) {
	if result, err := dynamo.ScanTable(&productTable); err != nil {
		return nil, err
	} else {
		var goods []model.Product
		for _, item := range result.Items {
			good := model.Product{}
			if err := dynamodbattribute.UnmarshalMap(item, &good); err != nil {
				return nil, err
			} else {
				goods = append(goods, good)
			}
		}
		return goods, nil
	}
}

func FindAllProductsByOwner(s *string) (*[]model.Product, error) {
	f := expression.Name("owner").Equal(expression.Value(s))
	if expr, err := expression.NewBuilder().WithFilter(f).Build(); err != nil {
		return nil, err
	} else {
		input := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			TableName:                 &productTable,
		}
		if result, err := dynamo.Scan(input); err != nil {
			return nil, err
		} else {
			var products []model.Product
			for _, item := range result.Items {
				product := model.Product{}
				if err := dynamodbattribute.UnmarshalMap(item, &product); err != nil {
					return nil, err
				} else {
					products = append(products, product)
				}
			}
			return &products, nil
		}
	}
}

func SaveProduct(p *model.Product) error {
	if p.Id == "" {
		if id, err := uuid.NewUUID(); err != nil {
			return err
		} else {
			p.Id = id.String()
		}
	}
	if item, err := dynamodbattribute.MarshalMap(&p); err != nil {
		return err
	} else {
		return dynamo.PutItem(item, &productTable)
	}
}
