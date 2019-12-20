package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	. "hc-api/service"
	"log"
	"os"
	"strings"
	"sync"
)

var orderTable = os.Getenv("ORDER_TABLE")
var packageTable = os.Getenv("PACKAGE_TABLE")
var en = expression.Name("status")
var c = expression.Or(
	expression.AttributeNotExists(en),
	expression.Equal(en, expression.Value("draft-1")),
	expression.Equal(en, expression.Value("draft-2")))

type Order struct {
	Id         string   `json:"id"`
	UserId     string   `json:"user_id"`
	AddressId  string   `json:"address_id_to"`
	PackageIds []string `json:"package_ids"`
	// User Profile
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	// User Address (where order is being shipped)
	Street string `json:"street,omitempty"`
	Unit   string `json:"unit,omitempty"`
	City   string `json:"city,omitempty"`
	State  string `json:"state,omitempty"`
	Zip5   string `json:"zip_5,omitempty"`
	Zip4   string `json:"zip_4,omitempty"`
	// report based data fields.
	OrderSum int64 `json:"order_sum,omitempty"`
	// transient variables, so to speak.
	Session   string `json:"session"`
	ProfileId string `json:"profile_id,omitempty"`
	// required for USPS
	Packages []Package `json:"packages"`
	// Package Id (ie Product Id) -> Vendor Id -> Service Id (description) -> Rate (price).
	Rates  map[string]map[string]map[string]string `json:"rates"`
	Vendor string                                  `json:"-"`
}

// Product information for order history data integrity.
// Package container dimensions.
// Transient variables.
// todo - use quantity to estimate shipping container dimensions
type Package struct {
	// ids
	Id            string `json:"id,omitempty"`
	ProductId     string `json:"product_id"`
	AddressIdFrom string `json:"address_id_from,omitempty"`
	AddressIdTo   string `json:"address_id_to,omitempty"`
	// product data
	ProductName  string `json:"product_name"`
	ProductPrice int64  `json:"product_price"`
	ProductQty   int    `json:"product_qty"`
	// usps
	ZipOrigination string `json:"zip_origination"`
	ZipDestination string `json:"zip_destination"`
	// fedex
	ShipperStateCode   string `json:"shipper_state_code"`
	RecipientStateCode string `json:"recipient_state_code"`
	// ups, fedex, usps
	ProductPounds int     `json:"pounds"`
	ProductOunces float32 `json:"ounces"`
	ProductWeight float32 `json:"product_weight"`
	ProductLength int     `json:"product_length"`
	ProductWidth  int     `json:"product_width"`
	ProductHeight int     `json:"product_height"`
	// ups, fedex
	TotalLength int     `json:"length"`
	TotalWidth  int     `json:"width"`
	TotalHeight int     `json:"height"`
	TotalWeight float32 `json:"weight"`
	// vendor data
	VendorName  string `json:"vendor_name"`
	VendorType  string `json:"vendor_type"`
	VendorPrice int    `json:"vendor_price"`
}

func NewOrder(body, ip string) (Order, error) {
	var o Order
	if err := json.Unmarshal([]byte(body), &o); err != nil {
		return o, err
	} else if userId, err := ValidateSession(o.Session, ip); err != nil {
		return o, err
	} else if id, err := uuid.NewUUID(); err != nil {
		return o, err
	} else {
		o.Id = id.String()
		o.UserId = userId
		o.PackageIds = make([]string, len(o.Packages))
		for i, p := range o.Packages {
			o.PackageIds = append(o.PackageIds, p.Id)
			p.TotalLength = p.ProductLength
			p.TotalWeight = p.ProductWeight * float32(p.ProductQty)
			p.TotalHeight = p.ProductHeight * p.ProductQty
			p.TotalWidth = p.ProductWidth * p.ProductQty

			aFrom, _ := base64.StdEncoding.DecodeString(p.AddressIdFrom)
			arrFrom := strings.Split(string(aFrom), ", ")
			p.ZipOrigination = strings.Split(arrFrom[len(arrFrom)-2], "-")[0]
			p.ShipperStateCode = arrFrom[len(arrFrom)-3]

			aTo, _ := base64.StdEncoding.DecodeString(o.AddressId)
			arrTo := strings.Split(string(aTo), ", ")
			p.ZipDestination = strings.Split(arrTo[len(arrTo)-2], "-")[0]
			p.RecipientStateCode = arrTo[len(arrTo)-3]

			o.Packages[i] = p
		}
		return o, nil
	}
}

func NewOrderResponse(body, ip, vendor string) (Order, error) {
	if o, err := NewOrder(body, ip); err != nil {
		return o, err
	} else {
		o.Vendor = vendor
		return o, err
	}
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	ip := r.RequestContext.Identity.SourceIP
	body := r.Body
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]\n", cmd, ip, body)

	switch cmd {

	case "calc-rates":
		if o, err := NewOrder(body, ip); err != nil {
			return BadGateway().Error(err).Build()
		} else {

			vs := []string{"UPS", "USPS"}

			orders := make(chan Order)

			var wg sync.WaitGroup
			wg.Add(len(vs))

			for _, v := range vs {
				go func(v string) {
					defer wg.Done()
					o, _ := NewOrderResponse(body, ip, v)
					o.Vendor = v
					err := Invoke().Handler("Shipping").QSP("cmd", "rate").QSP("v", v).Body(o).Marshal(&o)
					if err != nil {
						log.Println(err)
					}
					orders <- o
				}(v)
			}

			rates := map[string]map[string]map[string]string{}
			go func() {
				for o := range orders {
					for k, v := range o.Rates {
						if _, m := rates[k]; m {
							rates[k][o.Vendor] = v[o.Vendor]
						} else {
							rates[k] = v
						}
					}
				}
			}()

			wg.Wait()

			o.Rates = rates
			return Ok().Data(&o).Build()
		}

	case "save-order":
		if o, err := NewOrder(body, ip); err != nil {
			return BadGateway().Error(err).Build()
		} else {
			// save all packages
			for _, p := range o.Packages {
				if err := Put(p, &packageTable); err != nil {
					return BadRequest().Error(err).Build()
				} else {
					o.OrderSum += p.ProductPrice + int64(p.VendorPrice)
				}
			}
			// find pertinent data for saving order
			if um, err := Invoke().
				Handler("User").
				QSP("cmd", "find").
				QSP("session", o.Session).
				QSP("ip", ip).Build(); err != nil {
				return BadRequest().Error(err).Build()
			} else if pm, err := Invoke().
				Handler("UserProfile").
				QSP("cmd", "find").
				QSP("id", fmt.Sprintf("%v", um["profile_id"])).Build(); err != nil {
				return BadRequest().Error(err).Build()
			} else {
				// we dont need to do this here, but we will need to do it prior to submitting final order
				o.ProfileId = fmt.Sprintf("%v", um["profile_id"])
				o.Email = fmt.Sprintf("%v", pm["email"])
				o.Phone = fmt.Sprintf("%v", pm["phone"])
				o.FirstName = fmt.Sprintf("%v", pm["first_name"])
				o.LastName = fmt.Sprintf("%v", pm["last_name"])
				// clear irrelevant data
				o.Packages = nil
				o.Rates = nil
				if err := Put(o, &orderTable); err != nil {
					return InternalServerError().Error(err).Build()
				}
			}
			return Ok().Build()
		}

	default:
		return BadRequest().Data(r).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
