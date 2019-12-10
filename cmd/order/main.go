package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/service"
	"net/http"
	"os"
)

var orderTable = os.Getenv("ORDER_TABLE")
var packageTable = os.Getenv("PACKAGE_TABLE")

type Order struct {
	Id         string   `json:"id"`
	UserId     string   `json:"user_id"`
	Session    string   `json:"session"`
	PackageIds []string `json:"package_ids"`
	PriceInt   string   `json:"price_int"`
	PriceDec   string   `json:"price_dec"`
	AddressId  string   `json:"address_id"`
	Street     string   `json:"street"`
	Unit       string   `json:"unit,omitempty"`
	City       string   `json:"city"`
	State      string   `json:"state"`
	Zip5       string   `json:"zip_5"`
	Zip4       string   `json:"zip_4,omitempty"`
	Deleted    bool     `json:"deleted,omitempty"`
}

type Package struct {
	Id              string `json:"id,omitempty"`
	OrderId         string `json:"order_id"`
	ProductId       string `json:"product_id"`
	ProductName     string `json:"product_name"`
	ProductPriceInt string `json:"product_price_int"`
	ProductPriceDec string `json:"product_price_dec"`
	ProductQty      string `json:"product_qty"`
	ProductSum      string `json:"product_sum"` // total product cost
	ProductZip      string `json:"product_zip"` // zip
	ShipVendor      string `json:"ship_vendor"` // usps, ups, fedex
	ShipType        string `json:"ship_type"`   // priority, etc.
	ShipSum         string `json:"ship_sum"`    // total shipping costs
	PackageSum      string `json:"package_sum"` // total sum
}

func (o *Order) Validate() error {
	if b, err := json.Marshal(o); err != nil {
		return err
	} else if str, err := service.VerifyAddress(string(b)); err != nil {
		return err
	} else if err := json.Unmarshal([]byte(str), &o); err != nil {
		return err
	} else if o.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		o.Id = id.String()
		return nil
	}
}

func (p *Package) Validate() error {
	if p.ProductId == "" || p.ProductQty == "" || p.ProductPriceInt == "" || p.ProductPriceDec == "" {
		return fmt.Errorf("bad product information")
	} else if p.ProductZip == "" || p.ShipType == "" || p.ShipVendor == "" {
		return fmt.Errorf("bad shipping information")
	} else if p.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		p.Id = id.String()
		return nil
	}
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]", cmd, ip, body)

	switch cmd {

	case "save-order":
		var o Order
		if err := json.Unmarshal([]byte(body), &o); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := o.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := service.Put(o, &orderTable); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&o).Build()
		}

	case "save-package":
		var p Package
		if err := json.Unmarshal([]byte(body), &p); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := service.Put(p, &packageTable); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&p).Build()
		}

	case "update-order-packages":
		var u service.SliceUpdate
		if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := service.UpdateSlice(&u.Id, &u.Expression, &orderTable, &u.Val); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Build()
		}

	case "delete-order":
		id := r.QueryStringParameters["id"]
		if err := service.Delete(&id, &orderTable); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Build()
		}

	case "delete-package":
		id := r.QueryStringParameters["id"]
		if err := service.Delete(&id, &orderTable); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
