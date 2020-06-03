package ups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type NodeType struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}

type PostageRateRequest struct {
	RateRequest RateRequest `json:"RateRequest"`
	Packages    []Package   `json:"packages"`
}

type PostageRateResponse struct {
	RateResponse RateResponse `json:"RateResponse"`
	Packages     []Package    `json:"packages"`
}

type RateRequest struct {
	Shipment Shipment `json:"Shipment"`
}

type RateRequestResponse struct {
	RateResponse RateResponse `json:"RateResponse"`
}

type RateResponse struct {
	RatedShipment RatedShipment `json:"RatedShipment"`
}

type Response struct {
	ResponseStatus NodeType `json:"ResponseStatus"`
}

type RatedShipment struct {
	Service      NodeType     `json:"Service"`
	TotalCharges TotalCharges `json:"TotalCharges"`
}

type TotalCharges struct {
	Value string `json:"MonetaryValue"`
}

type Shipment struct {
	Shipper  ShippingEntity `json:"Shipper"`
	ShipTo   ShippingEntity `json:"ShipTo"`
	ShipFrom ShippingEntity `json:"ShipFrom"`
	Service  NodeType       `json:"Service"`
	Weight   WeightNode     `json:"ShipmentTotalWeight"`
	Package  Package        `json:"Package"`
}

type ShippingEntity struct {
	Address Address `json:"Address"`
}

type Address struct {
	PostalCode  string `json:"PostalCode"`
	CountryCode string `json:"CountryCode"`
}

type Package struct {
	Id         string `json:"id"`
	NodeType   `json:"PackagingType"`
	Dimensions `json:"Dimensions"`
	WeightNode `json:"PackageWeight"`
}

type Dimensions struct {
	DimensionType NodeType `json:"UnitOfMeasurement"`
	Length        string   `json:"Length"`
	Width         string   `json:"Width"`
	Height        string   `json:"Height"`
}

type WeightNode struct {
	WeightType NodeType `json:"UnitOfMeasurement"`
	Weight     string   `json:"Weight"`
}

const rateUrl = "https://onlinetools.ups.com/ship/v1/rating/Rate"

var (
	serviceMap = map[string]string{
		"01": "NEXT DAY AIR",
		"02": "2ND DAY AIR",
		"03": "GROUND",
		"11": "STANDARD",
		"12": "3-DAY SELECT",
		"14": "NEXT DAY AIR EARLY AM",
		"59": "NEXT DAY 2ND DAY AIR AM",
		"65": "UPS SAVER",
	}
	header = map[string][]string{
		"AccessLicenseNumber": {os.Getenv("UPS_ACCESS_LICENSE_NUMBER")},
		"Password":            {os.Getenv("UPS_PASSWORD")},
		"Content-Type":        {"application/json"},
		"transId":             {os.Getenv("UPS_TRANSACTION_ID")},
		"transSrc":            {os.Getenv("UPS_TRANSACTION_SRC")},
		"Username":            {os.Getenv("UPS_USERNAME")},
	}
)

func GetRates(in PostageRateRequest) (interface{}, error) {

	rates := map[string]map[string]map[string]string{}

	for _, p := range in.Packages {

		services := map[string]string{}

		for k, v := range serviceMap {

			in.RateRequest.Shipment.Package.NodeType = NodeType{k, v}

			out := PostageRateResponse{}

			if err := DoRequest(in, &out); err == nil {
				services[v] = out.RateResponse.RatedShipment.TotalCharges.Value
			}
		}

		if len(services) > 0 {
			rates[p.Id] = map[string]map[string]string{"UPS": services}
		}
	}

	return rates, nil
}

func DoRequest(in PostageRateRequest, out *PostageRateResponse) error {

	b, _ := json.Marshal(in)

	request, err := http.NewRequest(http.MethodPost, rateUrl, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 25 * time.Second}
	request.Header = header

	response, _ := client.Do(request)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(body, &out)
	_ = response.Body.Close()
	if out.RateResponse.RatedShipment.TotalCharges.Value == "" {
		return fmt.Errorf("received $0 service charge, service will be omitted from results")
	} else {
		return nil
	}
}
