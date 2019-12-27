package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	. "hc-api/service"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var t = os.Getenv("ORDER_TABLE")

type Order struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	OrderSum  int64  `json:"order_sum"`
	// Package Id (ie Product Id) -> Vendor Id -> Service Name -> Rate.
	Rates    map[string]map[string]map[string]string `json:"rates,omitempty"`
	Packages []Package                               `json:"packages"`
	// transient data (when outside of this layer)
	PackageIds []string `json:"package_ids,omitempty"`
	Session    string   `json:"session"`
	Vendor     string   `json:"-"`
	// auditing
	Created string `json:"created"`
}

// Product information for order history data integrity.
// Package container dimensions.
// Transient variables.
// todo - use quantity to estimate shipping container dimensions
type Package struct {
	// ids
	Id            string `json:"id,omitempty"`
	ProductId     string `json:"product_id"`
	AddressIdFrom string `json:"address_id_from"`
	AddressIdTo   string `json:"address_id_to"`
	// product data
	ProductName  string `json:"product_name"`
	ProductPrice int64  `json:"product_price"`
	ProductQty   int    `json:"product_qty"`
	ProductImg   string `json:"product_img,omitempty"`
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
	VendorPrice int64  `json:"vendor_price"`
	TotalPrice  int64  `json:"total_price"`
}

func NewOrder(body, ip string) (Order, error) {
	var o Order
	if err := json.Unmarshal([]byte(body), &o); err != nil {
		return o, err
	} else if userId, err := ValidateSession(o.Session, ip); err != nil {
		return o, err
	} else {
		if o.Id == "" {
			id, _ := uuid.NewUUID()
			o.Id = id.String()
		}
		o.UserId = userId
		o.PackageIds = make([]string, len(o.Packages))
		for i, p := range o.Packages {
			if p.Id == "" {
				continue
			}
			o.PackageIds = append(o.PackageIds, p.Id)
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
	fmt.Printf("REQUEST cmd=[%s] ip=[%s] body=[%s]\n", cmd, body, ip)

	switch cmd {

	case "find-by-ids":
		var p []Order
		ss := strings.Split(r.QueryStringParameters["ids"], ",")
		if err := FindAllById(t, ss, &p); err != nil {
			return BadRequest().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "calc-rates":
		var p Order
		if err := json.Unmarshal([]byte(body), &p); err != nil {
			return BadGateway().Error(err).Build()
		}
		fmt.Println(p)
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
			// sum all packages
			for _, p := range o.Packages {
				o.OrderSum += p.ProductPrice + p.VendorPrice
			}
			// find pertinent data for saving order
			if um, err := Invoke().Handler("User").IP(ip).CMD("find").Session(o.Session).Build(); err != nil {
				return BadRequest().Error(err).Build()
			} else if pm, err := Invoke().Handler("UserProfile").IP(ip).Session(o.Session).CMD("find").
				QSP("id", fmt.Sprintf("%v", um["profile_id"])).Build(); err != nil {
				return BadRequest().Error(err).Build()
			} else {
				// we dont need to do this here, but we will need to do it prior to submitting final order
				o.Email = fmt.Sprintf("%v", pm["email"])
				o.Phone = fmt.Sprintf("%v", pm["phone"])
				o.FirstName = fmt.Sprintf("%v", pm["first_name"])
				o.LastName = fmt.Sprintf("%v", pm["last_name"])
				o.PackageIds = nil
				o.Created = time.Now().UTC().Format(time.RFC3339)
				if err := Put(o, &t); err != nil {
					return InternalServerError().Error(err).Build()
				} else if _, err := Invoke().
					Handler("User").
					QSP("cmd", "update").
					Body(SliceUpdate{Session: o.Session, Val: []string{o.Id}, Expression: "add order_ids :p"}).
					Build(); err != nil {
					return BadRequest().Error(err).Build()
				} else {
					return Ok().Data(&o).Build()
				}
			}
		}

	default:
		return BadRequest().Data(r).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
