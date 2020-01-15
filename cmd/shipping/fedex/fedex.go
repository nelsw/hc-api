package fedex

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var serviceTypes = []string{
	//"PRIORITY_OVERNIGHT",
	//"FIRST_OVERNIGHT",
	//"STANDARD_OVERNIGHT",
	//"FEDEX_2_DAY_AM",
	//"FEDEX_2_DAY",
	//"FEDEX_EXPRESS_SAVER",
	"FEDEX_GROUND",
}

var pck = os.Getenv("FEDEX_PARENT_CREDENTIAL_KEY")
var pcp = os.Getenv("FEDEX_PARENT_CREDENTIAL_PASSWORD")
var uck = os.Getenv("FEDEX_USER_CREDENTIAL_KEY")
var ucp = os.Getenv("FEDEX_USER_CREDENTIAL_PASSWORD")
var can = os.Getenv("FEDEX_CLIENT_ACCOUNT_NUMBER")
var mtr = os.Getenv("FEDEX_METER_NUMBER")

type PackageRateRequest struct {
	Packages []PackageRequest                        `json:"packages"`
	Rates    map[string]map[string]map[string]string `json:"rates"`
}

type PackageRateResponse struct {
	Packages []PackageRequest                        `json:"packages"`
	Rates    map[string]map[string]map[string]string `json:"rates"`
}

type PackageRequest struct {
	Id                       string  `json:"id"`
	ParentCredentialKey      string  `json:"parent_credential_key"`
	ParentCredentialPassword string  `json:"parent_credential_password"`
	UserCredentialKey        string  `json:"user_credential_key"`
	UserCredentialPassword   string  `json:"user_credential_password"`
	ClientAccountNumber      string  `json:"client_account_number"`
	MeterNumber              string  `json:"meter_number"`
	ShipperStateCode         string  `json:"shipper_state_code"`
	ShipperZipCode           string  `json:"zip_origination"`
	ShipperCountryCode       string  `json:"shipper_country_code"`
	RecipientStateCode       string  `json:"recipient_state_code"`
	RecipientZipCode         string  `json:"zip_destination"`
	RecipientCountryCode     string  `json:"recipient_country_code"`
	ServiceType              string  `json:"service_type"`
	Width                    float32 `json:"width"`
	Length                   float32 `json:"length"`
	Height                   float32 `json:"height"`
	Weight                   float32 `json:"weight"`
	ShipTimestamp            string  `json:"ship_timestamp"`
}

func NewPackageRateRequest(s string) (PackageRateRequest, error) {
	var out PackageRateRequest
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return out, err
	} else {
		var packages []PackageRequest
		for i := 0; i < len(out.Packages); i++ {
			p := out.Packages[i]
			p.ParentCredentialKey = pck
			p.ParentCredentialPassword = pcp
			p.UserCredentialKey = uck
			p.UserCredentialPassword = ucp
			p.ClientAccountNumber = can
			p.MeterNumber = mtr
			p.RecipientCountryCode = "US"
			p.ShipperCountryCode = "US"

			//tt, err := time.Parse(time.RFC3339, time.Now().String())
			//if err != nil {
			//	panic(err)
			//}
			//p.ShipTimestamp = tt.String()
			for _, v := range serviceTypes {
				p.ServiceType = v
				packages = append(packages, p)
			}
		}
		out.Packages = packages
		return out, nil
	}
}

type RateResponse struct {
	XMLName      xml.Name     `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	ResponseBody ResponseBody `xml:"Entity"`

	Rates map[string]map[string]map[string]string `json:"rates"`
}

type ResponseBody struct {
	XMLName   xml.Name  `xml:"http://schemas.xmlsoap.org/soap/envelope/ Entity"`
	RateReply RateReply `xml:"RateReply"`
}

type RateReply struct {
	XMLName           xml.Name          `xml:"RateReply"`
	Notifications     Notifications     `xml:"Notifications"`
	TransactionDetail TransactionDetail `xml:"TransactionDetail"`
	RateReplyDetails  RateReplyDetails  `xml:"RateReplyDetails"`
}

type Notifications struct {
	XMLName  xml.Name `xml:"Notifications"`
	Severity string   `xml:"Severity"` // success, etc.
	Message  string   `xml:"Message"`
}

type TransactionDetail struct {
	XMLName               xml.Name `xml:"TransactionDetail"`
	CustomerTransactionId string   `xml:"CustomerTransactionId"`
}

type RateReplyDetails struct {
	XMLName              xml.Name             `xml:"RateReplyDetails"`
	ServiceType          string               `xml:"ServiceType"`
	RatedShipmentDetails RatedShipmentDetails `xml:"RatedShipmentDetails"`
}

type RatedShipmentDetails struct {
	XMLName            xml.Name           `xml:"RatedShipmentDetails"`
	ShipmentRateDetail ShipmentRateDetail `xml:"ShipmentRateDetail"`
}

type ShipmentRateDetail struct {
	XMLName    xml.Name   `xml:"ShipmentRateDetail"`
	GrandTotal GrandTotal `xml:"TotalNetChargeWithDutiesAndTaxes"`
}

type GrandTotal struct {
	XMLName xml.Name `xml:"TotalNetChargeWithDutiesAndTaxes"`
	Amount  float32  `xml:"Amount"`
}

const rateXml = `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns="http://fedex.com/ws/rate/v26">
   <SOAP-ENV:Entity>
      <RateRequest>
         <WebAuthenticationDetail>
			<ParentCredential>
               <Key>{{.ParentCredentialKey}}</Key>
               <Password>{{.ParentCredentialPassword}}</Password>
            </ParentCredential>
            <UserCredential>
               <Key>{{.UserCredentialKey}}</Key>
               <Password>{{.UserCredentialPassword}}</Password>
            </UserCredential>
         </WebAuthenticationDetail>
         <ClientDetail>
            <AccountNumber>{{.ClientAccountNumber}}</AccountNumber>
            <MeterNumber>{{.MeterNumber}}</MeterNumber>
         </ClientDetail>
         <Version>
            <ServiceId>crs</ServiceId>
            <Major>26</Major>
            <Intermediate>0</Intermediate>
            <Minor>0</Minor>
         </Version>
         <RequestedShipment>
            <ShipTimestamp>2019-07-24T12:34:56-06:00</ShipTimestamp>
            <DropoffType>REGULAR_PICKUP</DropoffType>
            <ServiceType>{{.ServiceType}}</ServiceType>
            <PackagingType>YOUR_PACKAGING</PackagingType>
            <TotalWeight>
               <Units>LB</Units>
               <Value>{{.Weight}}</Value>
            </TotalWeight>
           <Shipper>
               <Address>
                  <StateOrProvinceCode>{{.ShipperStateCode}}</StateOrProvinceCode>
                  <PostalCode>{{.ShipperZipCode}}</PostalCode>
                  <CountryCode>{{.ShipperCountryCode}}</CountryCode>
               </Address>
            </Shipper>
            <Recipient>
               <Address>
                  <StateOrProvinceCode>{{.RecipientStateCode}}</StateOrProvinceCode>
                  <PostalCode>{{.RecipientZipCode}}</PostalCode>
                  <CountryCode>{{.RecipientCountryCode}}</CountryCode>
               </Address>
            </Recipient>
            <ShippingChargesPayment>
               <PaymentType>SENDER</PaymentType>
            </ShippingChargesPayment>
            <RateRequestTypes>LIST</RateRequestTypes>
            <PackageCount>1</PackageCount>
            <RequestedPackageLineItems>
               <SequenceNumber>1</SequenceNumber>
               <GroupNumber>1</GroupNumber>
               <GroupPackageCount>1</GroupPackageCount>
               <Weight>
                  <Units>LB</Units>
                  <Value>{{.Weight}}</Value>
               </Weight>
               <Dimensions>
                  <Length>{{.Length}}</Length>
                  <Width>{{.Width}}</Width>
                  <Height>{{.Height}}</Height>
                  <Units>IN</Units>
               </Dimensions>
               <ContentRecords>
                  <PartNumber>123445</PartNumber>
                  <ItemNumber>kjdjalsro1262739827</ItemNumber>
                  <ReceivedQuantity>12</ReceivedQuantity>
                  <Description>ContentDescription</Description>
               </ContentRecords>
            </RequestedPackageLineItems>
         </RequestedShipment>
      </RateRequest>
   </SOAP-ENV:Entity>
</SOAP-ENV:Envelope>`

var callUrl = `https://ws.fedex.com:443/web-services` //os.Getenv("FEDEX_CALL_URL")
var soapTemplate *template.Template

func init() {
	t, err := template.New("InputRequest").Parse(rateXml)
	template.Must(t, err)
	soapTemplate = t
}

func GetRates(s string) (PackageRateResponse, error) {
	var out PackageRateResponse
	if in, err := NewPackageRateRequest(s); err != nil {
		return out, err
	} else {
		packages := map[string]map[string]map[string]string{}
		for _, p := range in.Packages {
			buffer := &bytes.Buffer{}
			encoder := xml.NewEncoder(buffer)
			client := &http.Client{Timeout: time.Duration(10 * time.Second)}
			doc := &bytes.Buffer{}
			if err := soapTemplate.Execute(doc, p); err != nil {
				return out, err
			} else if err = encoder.Encode(doc.String()); err != nil {
				return out, err
			} else if req, err := http.NewRequest(http.MethodPost, callUrl, bytes.NewBuffer([]byte(doc.String()))); err != nil {
				return out, err
			} else if resp, err := client.Do(req); err != nil {
				return out, err
			} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
				return out, err
			} else {
				defer resp.Body.Close()
				var r RateResponse
				fmt.Println(string(body))
				if err = xml.Unmarshal(body, &r); err != nil {
					return out, err
				} else if r.ResponseBody.RateReply.Notifications.Severity != "SUCCESS" {
					return out, fmt.Errorf(r.ResponseBody.RateReply.Notifications.Message)
				}
				a := r.ResponseBody.RateReply.RateReplyDetails.RatedShipmentDetails.ShipmentRateDetail.GrandTotal.Amount
				svc := map[string]string{p.ServiceType: strconv.Itoa(int(a))}
				ven := map[string]map[string]string{"FEDEX": svc}
				packages[p.Id] = ven
			}
		} // end packages loop
		out.Rates = packages
		return out, nil
	}
}
