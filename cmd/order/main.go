package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/service"
	"net/http"
	"os"
)

var orderTable = os.Getenv("ORDER_TABLE")
var packageTable = os.Getenv("PACKAGE_TABLE")
var c = "status IN (:s1, :s2)"
var e = map[string]*dynamodb.AttributeValue{":s1": {S: aws.String("draft-1")}, ":s2": {S: aws.String("draft-2")}}

type Order struct {
	Id         string    `json:"id,omitempty"`
	Session    string    `json:"session,omitempty"`
	UserId     string    `json:"user_id,omitempty"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Phone      string    `json:"phone"`
	AddressId  string    `json:"address_id,omitempty"`
	Street     string    `json:"street"`
	Unit       string    `json:"unit,omitempty"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Zip5       string    `json:"zip_5"`
	Zip4       string    `json:"zip_4,omitempty"`
	Status     string    `json:"status"` // draft-1, draft-1, draft-1, processing, delivered, complete.
	PackageSum int64     `json:"package_sum"`
	PostageSum int64     `json:"postage_sum"`
	OrderSum   int64     `json:"order_sum"`
	PackageIds []string  `json:"package_ids,omitempty"`
	Packages   []Package `json:"packages,omitempty"`
}

type Package struct {
	Id            string `json:"id,omitempty"`
	OrderId       string `json:"order_id"`
	ProductId     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	ProductPrice  int64  `json:"product_price"`
	ProductQty    int    `json:"product_qty"`
	ProductZip    string `json:"product_zip"`    // zip, for shipping rate calculations
	PostageVendor string `json:"postage_vendor"` // usps, ups, fedex
	PostageType   string `json:"postage_type"`   // priority, etc.
	PostagePrice  int64  `json:"postage_price"`
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
	if p.ProductId == "" || p.ProductQty < 1 || p.ProductPrice < 1 || p.PostagePrice < 0 {
		return fmt.Errorf("bad product information")
	} else if p.PostageType == "" || p.PostagePrice < 1 || p.PostageVendor == "" || p.ProductZip == "" {
		return fmt.Errorf("bad shipping information")
	} else if p.OrderId == "" {
		return fmt.Errorf("bad order information")
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
		} else if err := service.PutConditionally(o, &orderTable, &c, e); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&o).Build()
		}

	case "save-order-packages":
		var o Order
		if err := json.Unmarshal([]byte(body), &o); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			if o.Session != "" {
				if id, err := service.ValidateSession(o.Session, ip); err != nil {
					return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
				} else {
					o.UserId = id
				}
			} else {
				o.UserId = "guest"
			}
			for _, p := range o.Packages {
				p.OrderId = o.Id
				if err := p.Validate(); err != nil {
					return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
				} else if err := service.Put(p, &packageTable); err != nil {
					return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
				}
			}
			// todo - email confirmation
			return response.New().Code(http.StatusOK).Build()
		}

	case "update-order-package-ids":
		var u service.SliceUpdate
		if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if id, err := service.ValidateSession(u.Session, ip); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := service.UpdateSlice(&id, &u.Expression, &orderTable, &u.Val); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Build()
		}

	case "delete-package":
		id := r.QueryStringParameters["id"]
		if err := service.Delete(&id, &packageTable); err != nil {
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
