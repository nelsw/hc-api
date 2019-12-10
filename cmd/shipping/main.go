package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	response "github.com/nelsw/hc-util/aws"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	validateApi    = "http://production.shippingapis.com/ShippingAPI.dll?API=Verify&XML="
	rateRequestApi = "http://production.shippingapis.com/ShippingAPI.dll?API=RateV4&XML="
)

var uspsUserId = os.Getenv("USPS_USER_ID")

type AddressValidateRequest struct {
	XMLName  xml.Name `xml:"AddressValidateRequest"`
	UserId   string   `xml:"USERID,attr"`
	Revision string   `xml:"Revision"`
	Address  Address  `xml:"Address" json:"address"`
}

type AddressValidateResponse struct {
	XMLName xml.Name `xml:"AddressValidateResponse"`
	Address Address  `xml:"Address" json:"address"`
}

type Address struct {
	XMLName xml.Name `xml:"Address"`
	Id      string   `json:"id,omitempty" xml:"ID,attr"`
	Unit    string   `json:"unit,omitempty" xml:"Address1,omitempty"`
	Street  string   `json:"street,omitempty" xml:"Address2,omitempty"`
	City    string   `json:"city,omitempty" xml:"City"`
	State   string   `json:"state,omitempty" xml:"State"`
	Zip5    string   `json:"zip_5,omitempty" xml:"Zip5"`
	Zip4    string   `json:"zip_4,omitempty" xml:"Zip4"`
}

type OrderQuoteRequest struct {
	XMLName  xml.Name         `xml:"RateV4Request"`
	UserId   string           `xml:"USERID,attr"`
	Revision string           `xml:"Revision"`
	Packages []PackageRequest `json:"packages" xml:"Package"`
}

type OrderQuoteResponse struct {
	XMLName  xml.Name          `xml:"RateV4Response"`
	Packages []PackageResponse `json:"packages" xml:"Package"`
}

type PackageRequest struct {
	XMLName        xml.Name `xml:"Package"`
	Id             string   `json:"id" xml:"ID,attr"`
	Service        string   `json:"service" xml:"Service"`
	ZipOrigination string   `json:"zip_origination" xml:"ZipOrigination"`
	ZipDestination string   `json:"zip_destination" xml:"ZipDestination"`
	Pounds         float32  `json:"pounds" xml:"Pounds"`
	Ounces         float32  `json:"ounces" xml:"Ounces"`
	Container      string   `json:"container" xml:"Container"`
	Width          float32  `json:"width" xml:"Width"`
	Length         float32  `json:"length" xml:"Length"`
	Height         float32  `json:"height" xml:"Height"`
	Girth          float32  `json:"girth" xml:"Girth"`
	Machinable     string   `json:"machinable" xml:"Machinable"`
}

type Content struct {
	XMLName            xml.Name `xml:"Content"`
	ContentType        string   `xml:"ContentType" json:"content_type"`
	ContentDescription string   `xml:"ContentDescription" json:"content_description"`
}

type PackageResponse struct {
	XMLName        xml.Name `xml:"Package"`
	Id             string   `xml:"ID,attr" json:"id"`
	ZipOrigination string   `xml:"ZipOrigination" json:"zip_origination"`
	ZipDestination string   `xml:"ZipDestination" json:"zip_destination"`
	Pounds         int      `xml:"Pounds" json:"pounds"`
	Ounces         int      `xml:"Ounces" json:"ounces"`
	Postage        Postage  `xml:"Postage" json:"postage"`
}

type Postage struct {
	Id    string `xml:"CLASSID,attr" json:"id"`
	Price string `xml:"Rate" json:"price"`
	Date  string `xml:"CommitmentDate" json:"date"`
	Name  string `xml:"CommitmentName" json:"name"`
}

func (pr *PackageRequest) Validate() error {
	if pr.ZipOrigination == "" {
		return fmt.Errorf("bad zip origination ")
	} else if pr.ZipDestination == "" {
		return fmt.Errorf("bad zip destination")
	} else if pr.Pounds == 0 {
		return fmt.Errorf("bad pounds")
	} else if pr.Ounces == 0 {
		return fmt.Errorf("bad ounces")
	} else {
		pr.Service = "PRIORITY"
		pr.Container = "LG FLAT RATE BOX"
		pr.Machinable = "TRUE"
		return nil
	}
}

func (oqr *OrderQuoteRequest) Validate() error {
	for _, p := range oqr.Packages {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	oqr.Revision = "2"
	oqr.UserId = uspsUserId
	return nil
}

func (avr *AddressValidateRequest) Validate() error {
	avr.UserId = uspsUserId
	avr.Revision = "1"
	return nil
}

func (a *Address) String() string {
	var sb strings.Builder
	sb.WriteString(a.Street)
	sb.WriteString(", ")
	if a.Unit != "" {
		sb.WriteString(a.Unit)
		sb.WriteString(", ")
	}
	sb.WriteString(a.City)
	sb.WriteString(", ")
	sb.WriteString(a.State)
	sb.WriteString(", ")
	sb.WriteString(a.Zip5)
	sb.WriteString("-")
	sb.WriteString(a.Zip4)
	sb.WriteString(", ")
	sb.WriteString("United States")
	return sb.String()
}

func getXML(url string) (string, error) {
	if resp, err := http.Get(url); err != nil {
		return "", fmt.Errorf("GET error: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status error: %v", resp.StatusCode)
	} else if data, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", fmt.Errorf("read body: %v", err)
	} else {
		log.Println(string(data))
		return string(data), nil
	}
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	fmt.Printf("REQUEST [%s]: [%v]", cmd, r)

	switch cmd {

	case "verify":
		var in AddressValidateRequest
		var avr AddressValidateResponse
		if err := json.Unmarshal([]byte(r.Body), &in); err != nil {
			return response.New().Code(http.StatusBadGateway).Text(err.Error()).Build()
		} else if err := in.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if b, err := xml.Marshal(&in); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if s, err := getXML(validateApi + url.PathEscape(string(b))); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := xml.Unmarshal([]byte(s), &avr); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			avr.Address.Id = base64.StdEncoding.EncodeToString([]byte(avr.Address.String()))
			return response.New().Code(http.StatusOK).Data(&avr.Address).Build()
		}

	case "rate":
		var in OrderQuoteRequest
		var out OrderQuoteResponse
		if err := json.Unmarshal([]byte(r.Body), &in); err != nil {
			return response.New().Code(http.StatusBadGateway).Text(err.Error()).Build()
		} else if err := in.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if b, err := xml.Marshal(in); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if s, err := getXML(rateRequestApi + url.PathEscape(string(b))); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := xml.Unmarshal([]byte(s), &out); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&out.Packages).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
