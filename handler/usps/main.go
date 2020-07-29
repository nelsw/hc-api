package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/usps"
	"strings"
)

const (
	ValidateApi    = "http://production.shippingapis.com/ShippingAPI.dll?API=Verify&XML="
	RateRequestApi = "http://production.shippingapis.com/ShippingAPI.dll?API=RateV4&XML="
)

var uid = os.Getenv("USPS_USER_ID")

func getXML(url string) ([]byte, error) {
	resp, _ := http.Get(url)
	data, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return data, err
}

// USPS Handler can verify and validation (entity) or perform a rate request.
func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	apigwp.LogRequest(r)

	authResponse := client.Invoke("tokenHandler", events.APIGatewayProxyRequest{Path: "authenticate", Headers: r.Headers})
	if authResponse.StatusCode != 200 {
		return apigwp.HandleResponse(authResponse)
	}

	switch r.QueryStringParameters["path"] {

	case "validate":

		var a usps.Address
		_ = json.Unmarshal([]byte(r.Body), &a)

		inBytes, _ := xml.Marshal(&usps.AddressValidateRequest{uid, "1", a})

		outBytes, _ := getXML(ValidateApi + url.PathEscape(string(inBytes)))

		var out = usps.AddressValidateResponse{}
		_ = xml.Unmarshal(outBytes, &out)

		return apigwp.ProxyResponse(200, r.Headers, out.Address)

	case "rate":
		var pp []usps.PackageRequest
		_ = json.Unmarshal([]byte(r.Body), &pp)

		for i, p := range pp {
			p.Service = "PRIORITY"
			p.Container = "LG FLAT RATE BOX"
			p.Machinable = "TRUE"
			pp[i] = p
		}

		inBytes, _ := xml.Marshal(&usps.RateV4Request{uid, "2", pp})
		outBytes, _ := getXML(RateRequestApi + url.PathEscape(string(inBytes)))

		var out usps.RateV4Response
		_ = xml.Unmarshal(outBytes, &out)

		rates := map[string]map[string]map[string]string{}
		for _, p := range out.Packages {
			rates[p.Id] = map[string]map[string]string{"USPS": {strings.Split(p.Postage.Type, "&")[0]: p.Postage.Price}}
		}
		return apigwp.ProxyResponse(200, r.Headers, rates)
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
