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
	"net/http"
	"os"
	"strconv"
	"strings"
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
	Session        string    `json:"session,omitempty"`
	ProfileId      string    `json:"profile_id"`
	Packages       []Package `json:"packages,omitempty"`
	VendorPackages map[string]map[string]Package
}

// Product information for order history data integrity.
// Package container dimensions.
// Transient variables.
// todo - use quantity to estimate shipping container dimensions
type Package struct {
	Id            string `json:"id,omitempty"`
	ProductId     string `json:"product_id"`
	AddressIdFrom string `json:"address_id_from,omitempty"`

	ProductName   string  `json:"product_name"`
	ProductPrice  int64   `json:"product_price"`
	ProductWeight float32 `json:"product_weight"`
	ProductLength int     `json:"product_length"`
	ProductWidth  int     `json:"product_width"`
	ProductHeight int     `json:"product_height"`
	ProductQty    int     `json:"product_qty"`

	TotalLength int     `json:"length"`
	TotalWidth  int     `json:"width"`
	TotalHeight int     `json:"height"`
	TotalWeight float32 `json:"weight"`

	PostageVendor string  `json:"postage_vendor"`
	PostageType   string  `json:"postage_type"`
	PostagePrice  string  `json:"postage_price"`
	TotalPrice    float32 `json:"total_price"`

	// transient
	Postage     Postage `json:"postage"`
	AddressIdTo string  `json:"address_id_to,omitempty"`

	ZipOrigination string  `json:"zip_origination"`
	ZipDestination string  `json:"zip_destination"`
	ProductPounds  int     `json:"pounds"`
	ProductOunces  float32 `json:"ounces"`
}

type Postage struct {
	Id     string `json:"id"`
	Vendor string `json:"vendor"`       // usps, ups, fedex
	Type   string `json:"postage_type"` // priority, etc.
	Price  string `json:"postage_price"`
}

func (o *Order) Validate() error {
	if o.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		o.Id = id.String()
		return nil
	}
}

func (p *Package) Validate() error {
	if p.ProductId == "" || p.ProductQty < 1 {
		return fmt.Errorf("bad product information")
	} else if m, err := Invoke().Handler("Product").QSP("cmd", "find").QSP("id", p.ProductId).Build(); err != nil {
		return err
	} else if price, err := strconv.Atoi(fmt.Sprintf("%v", m["price"])); err != nil {
		return err
	} else if weight, err := strconv.Atoi(fmt.Sprintf("%v", m["weight"])); err != nil {
		return err
	} else if length, err := strconv.Atoi(fmt.Sprintf("%v", m["length"])); err != nil {
		return err
	} else if width, err := strconv.Atoi(fmt.Sprintf("%v", m["width"])); err != nil {
		return err
	} else if height, err := strconv.Atoi(fmt.Sprintf("%v", m["height"])); err != nil {
		return err
	} else {
		p.AddressIdFrom = fmt.Sprintf("%v", m["address_id"])
		p.ProductName = fmt.Sprintf("%v", m["name"])
		p.ProductWeight = float32(weight)
		p.ProductPounds = 1
		p.ProductOunces = 0
		p.ProductLength = int(length)
		p.ProductWidth = int(width)
		p.ProductHeight = int(height)
		p.ProductPrice = int64(price)
		p.TotalLength = p.ProductLength
		p.TotalWeight = p.ProductWeight * float32(p.ProductQty)
		p.TotalHeight = p.ProductHeight * p.ProductQty
		p.TotalWidth = p.ProductWidth * p.ProductQty
		bTo, _ := base64.StdEncoding.DecodeString(p.AddressIdFrom)
		aTo := strings.Split(string(bTo), ", ")
		p.ZipOrigination = strings.Split(aTo[len(aTo)-2], "-")[0]
		bFr, _ := base64.StdEncoding.DecodeString(p.AddressIdTo)
		aFr := strings.Split(string(bFr), ", ")
		p.ZipDestination = strings.Split(aFr[len(aFr)-2], "-")[0]
		if p.Id != "" {
			return nil
		} else if id, err := uuid.NewUUID(); err != nil {
			return err
		} else {
			p.Id = id.String()
			return nil
		}
	}
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	ip := r.RequestContext.Identity.SourceIP
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]\n", cmd, ip, r.Body)

	switch cmd {

	case "calc-rates":
		var o Order
		if err := json.Unmarshal([]byte(r.Body), &o); err != nil {
			return BadGateway().Error(err).Build()
		} else if userId, err := ValidateSession(o.Session, ip); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := o.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if um, err := Invoke().Handler("User").QSP("cmd", "find").QSP("session", o.Session).QSP("ip", ip).Build(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if pm, err := Invoke().Handler("UserProfile").QSP("cmd", "find").QSP("id", fmt.Sprintf("%v", um["profile_id"])).Build(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			o.ProfileId = fmt.Sprintf("%v", um["profile_id"])
			o.Status = cmd
			o.UserId = userId
			o.Email = fmt.Sprintf("%v", pm["email"])
			o.Phone = fmt.Sprintf("%v", pm["phone"])
			o.FirstName = fmt.Sprintf("%v", pm["first_name"])
			o.LastName = fmt.Sprintf("%v", pm["last_name"])
			for i, p := range o.Packages {
				p.AddressIdTo = o.AddressId
				if err := p.Validate(); err != nil {
					return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
				}
				o.Packages[i] = p
				if o.VendorPackages == nil {
					pm := map[string]Package{}
					pm[p.Id] = p
					vp := map[string]map[string]Package{}
					vp["USPS"] = pm
					o.VendorPackages = vp
				}
				pack := o.VendorPackages["USPS"][p.Id]
				pack.PostagePrice = p.PostagePrice
				o.VendorPackages["USPS"][p.Id] = pack
			}
			var newO Order
			if err := Invoke().Handler("Shipping").QSP("cmd", "rate").Body(o).Marshal(&newO); err != nil {
				return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
			} else {
				for _, p := range newO.Packages {
					pack := o.VendorPackages["USPS"][p.Id]
					pack.PostagePrice = p.PostagePrice
					o.VendorPackages["USPS"][p.Id] = pack
				}

				return Ok().Data(&o).Build()
			}

		}

	case "save-order":
		var o Order
		if err := json.Unmarshal([]byte(r.Body), &o); err != nil {
			return BadGateway().Error(err).Build()
		} else if err := o.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if exp, err := expression.NewBuilder().WithCondition(c).Build(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := PutConditionally(o, &orderTable, exp.Condition(), exp.Values()); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return Ok().Data(&o).Build()
		}

	case "save-order-packages":
		var o Order
		if err := json.Unmarshal([]byte(r.Body), &o); err != nil {
			return BadGateway().Error(err).Build()
		} else if err := o.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			if o.Session == "" {
				o.UserId = ip
			} else if id, err := ValidateSession(o.Session, ip); err != nil {
				return Unauthorized().Error(err).Build()
			} else {
				o.UserId = id
			}
			for _, p := range o.Packages {
				if err := p.Validate(); err != nil {
					return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
				} else if err := Put(p, &packageTable); err != nil {
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
