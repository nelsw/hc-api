package repo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"hc-api/model"
	"hc-api/repo/dynamo"
	"os"
)

var productTable = os.Getenv("PRODUCT_TABLE")

func FindAllProducts() ([]model.Product, error) {
	if result, err := dynamo.Scan(&productTable); err != nil {
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
