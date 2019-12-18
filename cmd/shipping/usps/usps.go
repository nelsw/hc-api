package usps

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	ValidateApi    = "http://production.shippingapis.com/ShippingAPI.dll?API=Verify&XML="
	RateRequestApi = "http://production.shippingapis.com/ShippingAPI.dll?API=RateV4&XML="
)

var uspsUserId = os.Getenv("USPS_USER_ID")

type VerificationRequest struct {
	XMLName  xml.Name `xml:"AddressValidateRequest"`
	UserId   string   `xml:"USERID,attr"`
	Revision string   `xml:"Revision"`
	Address  Address  `xml:"Address" json:"address"`
}

type VerificationResponse struct {
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

type RateRequest struct {
	XMLName  xml.Name                                `xml:"RateV4Request"`
	UserId   string                                  `xml:"USERID,attr"`
	Revision string                                  `xml:"Revision"`
	Packages []PackageRequest                        `xml:"Package" json:"packages"`
	Rates    map[string]map[string]map[string]string `json:"rates" xml:"-"`
}

type RateResponse struct {
	XMLName  xml.Name                                `xml:"RateV4Response"`
	Packages []PackageResponse                       `json:"packages" xml:"Package"`
	Rates    map[string]map[string]map[string]string `json:"rates"`
}

type PackageRequest struct {
	XMLName        xml.Name `xml:"Package"`
	Id             string   `json:"id" xml:"ID,attr"`
	Service        string   `json:"service,omitempty" xml:"Service"`
	ZipOrigination string   `json:"zip_origination" xml:"ZipOrigination"`
	ZipDestination string   `json:"zip_destination" xml:"ZipDestination"`
	Pounds         int      `json:"pounds" xml:"Pounds"`
	Ounces         float32  `json:"ounces" xml:"Ounces"`
	Container      string   `json:"container" xml:"Container"`
	Width          float32  `json:"width" xml:"Width"`
	Length         float32  `json:"length" xml:"Length"`
	Height         float32  `json:"height" xml:"Height"`
	Girth          float32  `json:"girth" xml:"Girth"`
	Machinable     string   `json:"machinable" xml:"Machinable"`
	AddressIdTo    string   `json:"address_id_to" xml:"-"`
	AddressIdFrom  string   `json:"address_id_from" xml:"-"`
	Vendor         string   `json:"vendor" xml:"-"`
}

type PackageResponse struct {
	XMLName        xml.Name `xml:"Package"`
	Id             string   `xml:"ID,attr" json:"id"`
	ZipOrigination string   `xml:"ZipOrigination" json:"zip_origination"`
	ZipDestination string   `xml:"ZipDestination" json:"zip_destination"`
	Pounds         int      `xml:"Pounds" json:"pounds"`
	Ounces         float32  `xml:"Ounces" json:"ounces"`
	Postage        Postage  `xml:"Postage" json:"postage"`
	Type           string   `xml:"-" json:"postage_type"`
	Price          string   `xml:"-" json:"postage_price"`
	Vendor         string   `xml:"-" json:"postage_vendor"`
}

type Postage struct {
	XMLName xml.Name `xml:"Postage"`
	Id      string   `xml:"CLASSID,attr" json:"id"`
	Type    string   `xml:"MailService" json:"type"`
	Price   string   `xml:"Rate" json:"price"`
}

func (p *PackageRequest) Validate() error {
	if p.ZipOrigination == "" {
		return fmt.Errorf("bad zip origination ")
	} else if p.ZipDestination == "" {
		return fmt.Errorf("bad zip destination")
	} else if p.Pounds == 0 {
		return fmt.Errorf("bad pounds")
	} else {
		return nil
	}
}

func (r *RateRequest) Validate() error {
	for i, p := range r.Packages {
		if err := p.Validate(); err != nil {
			return err
		}
		p.Service = "PRIORITY"
		p.Container = "LG FLAT RATE BOX"
		p.Machinable = "TRUE"
		r.Packages[i] = p
	}
	r.Revision = "2"
	r.UserId = uspsUserId
	return nil
}

func (avr *VerificationRequest) Validate() error {
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

func GetAddress(s string) (Address, error) {
	var in VerificationRequest
	var out VerificationResponse
	if err := json.Unmarshal([]byte(s), &in); err != nil {
		return out.Address, err
	} else if err := in.Validate(); err != nil {
		return out.Address, err
	} else if b, err := xml.Marshal(&in); err != nil {
		return out.Address, err
	} else if s, err := getXML(ValidateApi + url.PathEscape(string(b))); err != nil {
		return out.Address, err
	} else if err := xml.Unmarshal([]byte(s), &out); err != nil {
		return out.Address, err
	} else {
		out.Address.Id = base64.StdEncoding.EncodeToString([]byte(out.Address.String()))
		return out.Address, nil
	}
}

func GetPostage(s string) (RateResponse, error) {
	var in RateRequest
	var out RateResponse
	if err := json.Unmarshal([]byte(s), &in); err != nil {
		return out, err
	} else if err := in.Validate(); err != nil {
		return out, err
	} else if b, err := xml.Marshal(in); err != nil {
		return out, err
	} else if s, err := getXML(RateRequestApi + url.PathEscape(string(b))); err != nil {
		return out, err
	} else if err := xml.Unmarshal([]byte(s), &out); err != nil {
		return out, err
	} else {
		packages := map[string]map[string]map[string]string{}
		for _, p := range out.Packages {
			service := map[string]string{}
			// for each service...
			t := strings.Split(p.Postage.Type, "&")[0]
			service[t] = p.Postage.Price
			// end for each service.
			vendor := map[string]map[string]string{}
			vendor["USPS"] = service
			packages[p.Id] = vendor
		}
		out.Rates = packages
		return out, nil
	}
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