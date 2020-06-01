package main

import (
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sam-app/pkg/model/usps"
	"strings"
)

const (
	ValidateApi    = "http://production.shippingapis.com/ShippingAPI.dll?API=Verify&XML="
	RateRequestApi = "http://production.shippingapis.com/ShippingAPI.dll?API=RateV4&XML="
)

var (
	uid     = os.Getenv("USPS_USER_ID")
	ErrOp   = fmt.Errorf("bad request\n")
	ErrHttp = fmt.Errorf("ERROR - unsuccessful http.Get\n")
	ErrIo   = fmt.Errorf("ERROR - unable to read response body\n")
)

// USPS Handler can verify and validation (entity) or perform a rate request.
func Handle(request usps.Request) (interface{}, error) {

	if request.Op == "validate" {
		v := usps.AddressValidateRequest{uid, "1", request.Address}

		inBytes, _ := xml.Marshal(&v)

		outBytes, _ := getXML(ValidateApi + url.PathEscape(string(inBytes)))

		var out = usps.AddressValidateResponse{}
		_ = xml.Unmarshal(outBytes, &out)

		return out.Address, nil
	}

	if request.Op == "rate" {
		for n, p := range request.Packages {

			p.Service = "PRIORITY"
			p.Container = "LG FLAT RATE BOX"
			p.Machinable = "TRUE"
			request.Packages[n] = p
		}

		in := usps.RateV4Request{uid, "2", request.Packages}

		inBytes, _ := xml.Marshal(&in)
		outBytes, _ := getXML(RateRequestApi + url.PathEscape(string(inBytes)))

		var out usps.RateV4Response
		_ = xml.Unmarshal(outBytes, &out)

		var k, v string
		rates := map[string]map[string]map[string]string{}
		for _, p := range out.Packages {
			k = strings.Split(p.Postage.Type, "&")[0]
			v = p.Postage.Price
			rates[p.Id] = map[string]map[string]string{"USPS": {k: v}}
		}

		return rates, nil
	}

	return nil, ErrOp
}

func getXML(url string) ([]byte, error) {
	if resp, err := http.Get(url); err != nil || resp.StatusCode != 200 {
		return nil, ErrHttp
	} else if data, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, ErrIo
	} else {
		_ = resp.Body.Close()
		return data, nil
	}
}

func main() {
	lambda.Start(Handle)
}
