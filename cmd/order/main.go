package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	response "github.com/nelsw/hc-util/aws"
	. "hc-api/service"
	"log"
	"net/http"
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
	Id         string   `json:"id,omitempty"`
	UserId     string   `json:"user_id,omitempty"`
	AddressId  string   `json:"address_id_to,omitempty"`
	PackageIds []string `json:"package_ids,omitempty"`
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
	// draft-1, draft-2, draft-3, processing, delivered, complete.
	Status string `json:"status,omitempty"`
	// report based data fields.
	OrderSum int64 `json:"order_sum,omitempty"`
	// transient variables, so to speak.
	Session   string `json:"session"`
	ProfileId string `json:"profile_id,omitempty"`
	// required for USPS
	Packages []Package `json:"packages,omitempty"`
	// Package Id (ie Product Id) -> Vendor Id -> Service Id (description) -> Rate (price).
	Rates  map[string]map[string]map[string]string `json:"rates,omitempty"`
	Vendor string                                  `json:"-"`
}

// Product information for order history data integrity.
// Package container dimensions.
// Transient variables.
// todo - use quantity to estimate shipping container dimensions
type Package struct {
	Id            string `json:"id,omitempty"`
	ProductId     string `json:"product_id"`
	AddressIdFrom string `json:"address_id_from,omitempty"`

	ZipOrigination string `json:"zip_origination"`
	ZipDestination string `json:"zip_destination"`

	ShipperStateCode   string `json:"shipper_state_code"`
	RecipientStateCode string `json:"recipient_state_code"`

	ProductName  string `json:"product_name"`
	ProductPrice int64  `json:"product_price"`

	ProductPounds int     `json:"pounds"`
	ProductOunces float32 `json:"ounces"`
	ProductWeight float32 `json:"product_weight"`

	ProductLength int `json:"product_length"`
	ProductWidth  int `json:"product_width"`
	ProductHeight int `json:"product_height"`

	ProductQty int `json:"product_qty"`

	TotalLength int     `json:"length"`
	TotalWidth  int     `json:"width"`
	TotalHeight int     `json:"height"`
	TotalWeight float32 `json:"weight"`

	TotalPrice float32 `json:"total_price"`

	// transient
	AddressIdTo string `json:"address_id_to,omitempty"`
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
		for i, p := range o.Packages {
			p.Id = p.ProductId
			p.AddressIdTo = o.AddressId
			p.TotalLength = p.ProductLength
			p.TotalWeight = p.ProductWeight * float32(p.ProductQty)
			p.TotalHeight = p.ProductHeight * p.ProductQty
			p.TotalWidth = p.ProductWidth * p.ProductQty

			aFrom, _ := base64.StdEncoding.DecodeString(p.AddressIdFrom)
			arrFrom := strings.Split(string(aFrom), ", ")
			p.ZipOrigination = strings.Split(arrFrom[len(arrFrom)-2], "-")[0]
			p.ShipperStateCode = arrFrom[len(arrFrom)-3]

			aTo, _ := base64.StdEncoding.DecodeString(p.AddressIdTo)
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
			return BadRequest().Error(err).Build()
		} else if um, err := Invoke().Handler("User").QSP("cmd", "find").QSP("session", o.Session).QSP("ip", ip).Build(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if pm, err := Invoke().Handler("UserProfile").QSP("cmd", "find").QSP("id", fmt.Sprintf("%v", um["profile_id"])).Build(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			// we dont need to do this here, but we will need to do it prior to submitting final order
			o.ProfileId = fmt.Sprintf("%v", um["profile_id"])
			o.Status = cmd
			o.Email = fmt.Sprintf("%v", pm["email"])
			o.Phone = fmt.Sprintf("%v", pm["phone"])
			o.FirstName = fmt.Sprintf("%v", pm["first_name"])
			o.LastName = fmt.Sprintf("%v", pm["last_name"])

			vs := []string{"UPS", "USPS", "FEDEX"}

			orders := make(chan Order)

			var wg sync.WaitGroup
			wg.Add(len(vs))

			for _, v := range vs {
				go func(v string) {
					defer wg.Done()
					o, _ := NewOrderResponse(body, ip, v)
					err := Invoke().Handler("Shipping").QSP("cmd", "rate").QSP("v", v).Body(o).Marshal(&o)
					if err != nil {
						log.Fatal(err)
					} else {
						orders <- o
					}
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
			return BadRequest().Error(err).Build()
		} else if exp, err := expression.NewBuilder().WithCondition(c).Build(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := PutConditionally(o, &orderTable, exp.Condition(), exp.Values()); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return Ok().Data(&o).Build()
		}

	case "save-order-packages":
		if o, err := NewOrder(body, ip); err != nil {
			return BadRequest().Error(err).Build()
		} else {
			if o.Session == "" {
				o.UserId = ip
			} else if id, err := ValidateSession(o.Session, ip); err != nil {
				return Unauthorized().Error(err).Build()
			} else {
				o.UserId = id
			}
			for _, p := range o.Packages {
				if err := Put(p, &packageTable); err != nil {
					return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
				} else {
					o.PackageIds = append(o.PackageIds, p.Id)
				}
			}
			if exp, err := expression.NewBuilder().WithCondition(c).Build(); err != nil {
				return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
			} else if err := PutConditionally(o, &orderTable, exp.Condition(), exp.Values()); err != nil {
				return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
			} else {
				// todo - email confirmation
				return Ok().Build()
			}
		}

	case "update-order-package-ids":
		var u SliceUpdate
		if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if id, err := ValidateSession(u.Session, ip); err != nil {
			return Unauthorized().Error(err).Build()
		} else if err := UpdateSlice(&id, &u.Expression, &orderTable, &u.Val); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Build()
		}

	default:
		return BadRequest().Data(r).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
