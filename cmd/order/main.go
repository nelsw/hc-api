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
	Id         string   `json:"id,omitempty"`
	Session    string   `json:"session,omitempty"`
	UserId     string   `json:"user_id,omitempty"`
	Email      string   `json:"email"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Phone      string   `json:"phone"`
	AddressId  string   `json:"address_id,omitempty"`
	Street     string   `json:"street"`
	Unit       string   `json:"unit,omitempty"`
	City       string   `json:"city"`
	State      string   `json:"state"`
	Zip5       string   `json:"zip_5"`
	Zip4       string   `json:"zip_4,omitempty"`
	Status     string   `json:"status"` // draft-1, draft-1, draft-1, processing, delivered, complete.
	PackageSum int64    `json:"package_sum"`
	PostageSum int64    `json:"postage_sum"`
	OrderSum   int64    `json:"order_sum"`
	PackageIds []string `json:"package_ids,omitempty"`
}

type Package struct {
	Id            string `json:"id,omitempty"`
	OrderId       string `json:"order_id"`
	ProductId     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	ProductSum    int64  `json:"product_price_dec"`
	ProductQty    string `json:"product_qty"`
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
	if p.ProductId == "" || p.ProductQty == "" || p.ProductPriceInt == "" || p.ProductPriceDec == "" {
		return fmt.Errorf("bad product information")
	} else if p.PostageType == "" || p.PostagePrice == "" || p.PostageVendor == "" || p.ProductZip == "" {
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
		} else if err := service.Put(o, &orderTable); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if token, err := service.NewSession(o.Id, ip); err != nil {
			return response.New().Code(http.StatusInternalServerError).Build()
		} else {
			o.Session = token
			return response.New().Code(http.StatusOK).Toke(token).Data(&o).Build()
		}

	case "save-order-packages":
		var pp []Package
		if err := json.Unmarshal([]byte(body), &pp); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			// todo - go routine
			for i, p := range pp {
				if err := p.Validate(); err != nil {
					return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
				} else if err := service.Put(p, &packageTable); err != nil {
					return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
				} else {
					pp[i] = p
				}
			}
			return response.New().Code(http.StatusOK).Data(&pp).Build()
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
