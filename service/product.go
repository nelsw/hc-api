package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"hc-api/model"
	"hc-api/service/dynamo"
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

func FindAllProductsByIds(ss *[]string) (*[]model.Product, error) {
	var pp []model.Product
	var keys []map[string]*dynamodb.AttributeValue
	for _, s := range *ss {
		keys = append(keys, map[string]*dynamodb.AttributeValue{"id": {S: aws.String(s)}})
	}
	if results, err := dynamo.GetBatch(keys, productTable); err != nil {
		return nil, err
	} else if err := dynamodbattribute.UnmarshalListOfMaps(results.Responses[productTable], &pp); err != nil {
		return nil, err
	} else {
		return &pp, nil
	}
}

func SaveProduct(p *model.Product) error {
	if item, err := dynamodbattribute.MarshalMap(&p); err != nil {
		return err
	} else {
		return dynamo.PutItem(item, &productTable)
	}
}
