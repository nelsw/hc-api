package ups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Type struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}

type PostageRateRequest struct {
	RateRequest RateRequest                             `json:"RateRequest"`
	Packages    []Package                               `json:"packages"`
	Rates       map[string]map[string]map[string]string `json:"rates"`
}

type PostageRateResponse struct {
	RateResponse RateResponse                            `json:"RateResponse"`
	Packages     []Package                               `json:"packages"`
	Rates        map[string]map[string]map[string]string `json:"rates"`
}

type RateRequest struct {
	Shipment Shipment `json:"Shipment"`
}

type RateRequestResponse struct {
	RateResponse RateResponse `json:"RateResponse"`
}

type RateResponse struct {
	Response      Response      `json:"Response"`
	RatedShipment RatedShipment `json:"RatedShipment"`
}

type Response struct {
	ResponseStatus Type `json:"ResponseStatus"`
}

// package - response
type RatedShipment struct {
	Service      Type         `json:"Service"`
	TotalCharges TotalCharges `json:"TotalCharges"`
}

type TotalCharges struct {
	Value string `json:"MonetaryValue"`
}

// Package
type Shipment struct {
	Shipper  ShippingEntity      `json:"Shipper"`
	ShipTo   ShippingEntity      `json:"ShipTo"`
	ShipFrom ShippingEntity      `json:"ShipFrom"`
	Service  Type                `json:"Service"`
	Weight   ShipmentTotalWeight `json:"ShipmentTotalWeight"`
	Package  Package             `json:"Package"`
}

type ShippingEntity struct {
	Address Address `json:"Address"`
}

type Address struct {
	PostalCode  string `json:"PostalCode"`
	CountryCode string `json:"CountryCode"`
}

type ShipmentTotalWeight struct {
	Type   Type   `json:"UnitOfMeasurement"`
	Weight string `json:"weight"`
}

type Package struct {
	Id            string `json:"id"`
	Type          `json:"PackagingType"`
	Dimensions    `json:"Dimensions"`
	PackageWeight `json:"PackageWeight"`

	ZipOrigination string `json:"zip_origination"`
	ZipDestination string `json:"zip_destination"`

	Weight float32 `json:"product_weight"`
	Length int     `json:"product_length"`
	Width  int     `json:"product_width"`
	Height int     `json:"product_height"`
}

type Dimensions struct {
	Type   Type   `json:"UnitOfMeasurement"`
	Length string `json:"Length"`
	Width  string `json:"Width"`
	Height string `json:"Height"`
}

type PackageWeight struct {
	Type   Type   `json:"UnitOfMeasurement"`
	Weight string `json:"Weight"`
}

var header = map[string][]string{
	"AccessLicenseNumber": {os.Getenv("UPS_ACCESS_LICENSE_NUMBER")},
	"Password":            {os.Getenv("UPS_PASSWORD")},
	"Content-Type":        {"application/json"},
	"transId":             {os.Getenv("UPS_TRANSACTION_ID")},
	"transSrc":            {os.Getenv("UPS_TRANSACTION_SRC")},
	"Username":            {os.Getenv("UPS_USERNAME")},
}

const rateUrl = "https://onlinetools.ups.com/ship/v1/rating/Rate"

var serviceTypeMap = map[string]string{
	"01": "NEXT DAY AIR",
	"02": "2ND DAY AIR",
	"03": "GROUND",
	"11": "STANDARD",
	"12": "3-DAY SELECT",
	"14": "NEXT DAY AIR EARLY AM",
	"59": "NEXT DAY 2ND DAY AIR AM",
	"65": "UPS SAVER",
}

func GetRates(s string) (PostageRateResponse, error) {
	var out PostageRateResponse
	var prr PostageRateRequest
	if err := json.Unmarshal([]byte(s), &prr); err != nil {
		return out, err
	} else {
		packages := map[string]map[string]map[string]string{}
		for _, p := range prr.Packages {
			fmt.Println(p)
			service := map[string]string{}
			for k, v := range serviceTypeMap {
				shipment := NewShipment(p.ZipOrigination, p.ZipDestination, p.Weight, p.Length, p.Width, p.Height, Type{k, v})
				prr.RateRequest.Shipment = shipment
				if o, err := DoRequest(prr); err != nil {
					return o, err
				} else if o.RateResponse.RatedShipment.TotalCharges.Value != "" {
					service[v] = o.RateResponse.RatedShipment.TotalCharges.Value
					fmt.Println(o)
				}

			}
			packages[p.Id] = map[string]map[string]string{"UPS": service}
		}
		out.Rates = packages
		return out, nil
	}
}

func DoRequest(prr PostageRateRequest) (PostageRateResponse, error) {
	var out PostageRateResponse
	b, err := json.Marshal(prr)
	if err != nil {
		return out, err
	}
	request, err := http.NewRequest(http.MethodPost, rateUrl, bytes.NewBuffer(b))
	if err != nil {
		return out, err
	} else {
		fmt.Println(prr)
		fmt.Println(&prr)
		client := &http.Client{Timeout: time.Duration(5 * time.Second)}
		request.Header = header
		if response, err := client.Do(request); err != nil {
			return out, err
		} else if body, err := ioutil.ReadAll(response.Body); err != nil {
			return out, err
		} else if err := json.Unmarshal(body, &out); err != nil {
			return out, err
		} else {
			fmt.Println(string(body))
			return out, nil
		}
	}
}

// think of this as a package
func NewShipment(s1, s2 string, s3 float32, s4, s5, s6 int, s0 Type) Shipment {
	return Shipment{
		Shipper:  ShippingEntity{Address: Address{PostalCode: s1, CountryCode: "US"}},
		ShipTo:   ShippingEntity{Address: Address{PostalCode: s2, CountryCode: "US"}},
		ShipFrom: ShippingEntity{Address: Address{PostalCode: s1, CountryCode: "US"}},
		Service:  s0,
		Weight:   ShipmentTotalWeight{Type: Type{Code: "LBS", Description: "Pounds"}, Weight: strconv.Itoa(int(s3))},
		Package: Package{
			Type: s0,
			Dimensions: Dimensions{
				Type{Code: "IN"},
				strconv.Itoa(s4),
				strconv.Itoa(s5),
				strconv.Itoa(s6),
			},
			PackageWeight: PackageWeight{Type: Type{Code: "LBS", Description: "Pounds"}, Weight: strconv.Itoa(int(s3))},
		},
	}
}
