package fedex

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type RateRequest struct {
	RecipientStateCode string    `json:"recipient_state_code"`
	RecipientZipCode   string    `json:"zip_destination"`
	Packages           []Package `json:"packages"`
}

type Package struct {
	Id string `json:"id"`

	ShipperStateCode string `json:"shipper_state_code"`
	ShipperZipCode   string `json:"zip_origination"`

	Width  float32 `json:"width"`
	Length float32 `json:"length"`
	Height float32 `json:"height"`
	Weight float32 `json:"weight"`

	RecipientStateCode, RecipientZipCode, ServiceType string
}

type RateResponse struct {
	XMLName      xml.Name     `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	ResponseBody ResponseBody `xml:"Body"`
}

type ResponseBody struct {
	XMLName   xml.Name  `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	RateReply RateReply `xml:"RateReply"`
}

type RateReply struct {
	Notifications     Notifications     `xml:"Notifications"`
	TransactionDetail TransactionDetail `xml:"TransactionDetail"`
	RateReplyDetails  RateReplyDetails  `xml:"RateReplyDetails"`
}

type Notifications struct {
	Severity string `xml:"Severity"` // must be success
	Message  string `xml:"Message"`
}

type TransactionDetail struct {
	CustomerTransactionId string `xml:"CustomerTransactionId"`
}

type RateReplyDetails struct {
	ServiceType          string               `xml:"ServiceType"`
	RatedShipmentDetails RatedShipmentDetails `xml:"RatedShipmentDetails"`
}

type RatedShipmentDetails struct {
	ShipmentRateDetail ShipmentRateDetail `xml:"ShipmentRateDetail"`
}

type ShipmentRateDetail struct {
	GrandTotal TotalNetChargeWithDutiesAndTaxes `xml:"TotalNetChargeWithDutiesAndTaxes"`
}

type TotalNetChargeWithDutiesAndTaxes struct {
	Amount float32 `xml:"Amount"`
}

var (
	soapTemplate *template.Template
	serviceTypes = []string{
		"PRIORITY_OVERNIGHT",
		"FIRST_OVERNIGHT",
		"STANDARD_OVERNIGHT",
		"FEDEX_2_DAY_AM",
		"FEDEX_2_DAY",
		"FEDEX_EXPRESS_SAVER",
		"FEDEX_GROUND",
	}
	url = os.Getenv("FEDEX_WS_URL")
)

func init() {
	var rateXml = `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns="http://fedex.com/ws/rate/v26">
	<SOAP-ENV:Body>
  	<RateRequest>
    	<WebAuthenticationDetail>
				<ParentCredential>
        	<Key>{{.pck}}</Key>
					<Password>{{.pcp}}</Password>
				</ParentCredential>
				<UserCredential>
				 <Key>{{.uck}}</Key>
				 <Password>{{.ucp}}</Password>
				</UserCredential>
			</WebAuthenticationDetail>
      <ClientDetail>
      	<AccountNumber>{{.can}}</AccountNumber>
        <MeterNumber>{{.mn}}</MeterNumber>
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
            <CountryCode>US</CountryCode>
					</Address>
				</Shipper>
        <Recipient>
        	<Address>
          	<StateOrProvinceCode>{{.RecipientStateCode}}</StateOrProvinceCode>
            <PostalCode>{{.RecipientZipCode}}</PostalCode>
            <CountryCode>US</CountryCode>
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
 </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
	rateXml = strings.ReplaceAll(rateXml, `{{.pck}}`, os.Getenv("FEDEX_PARENT_CREDENTIAL_KEY"))
	rateXml = strings.ReplaceAll(rateXml, `{{.pcp}}`, os.Getenv("FEDEX_PARENT_CREDENTIAL_PASSWORD"))
	rateXml = strings.ReplaceAll(rateXml, `{{.uck}}`, os.Getenv("FEDEX_USER_CREDENTIAL_KEY"))
	rateXml = strings.ReplaceAll(rateXml, `{{.ucp}}`, os.Getenv("FEDEX_USER_CREDENTIAL_PASSWORD"))
	rateXml = strings.ReplaceAll(rateXml, `{{.can}}`, os.Getenv("FEDEX_CLIENT_ACCOUNT_NUMBER"))
	rateXml = strings.ReplaceAll(rateXml, `{{.mn}}`, os.Getenv("FEDEX_METER_NUMBER"))
	t, err := template.New("xml").Parse(rateXml)
	template.Must(t, err)
	soapTemplate = t
}

func GetRates(in RateRequest) (interface{}, error) {

	rates := map[string]map[string]map[string]string{}

	for _, p := range in.Packages {

		services := map[string]string{}

		p.RecipientZipCode = in.RecipientZipCode
		p.RecipientStateCode = in.RecipientStateCode

		for _, s := range serviceTypes {

			p.ServiceType = s

			doc := &bytes.Buffer{}
			_ = soapTemplate.Execute(doc, p)

			req, _ := http.NewRequest(http.MethodPost, url, doc)

			client := &http.Client{Timeout: 25 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			_ = resp.Body.Close()

			var r RateResponse
			_ = xml.Unmarshal(body, &r)

			if r.ResponseBody.RateReply.Notifications.Severity != "SUCCESS" {
				return nil, fmt.Errorf(r.ResponseBody.RateReply.Notifications.Message)
			}

			a := r.ResponseBody.RateReply.RateReplyDetails.RatedShipmentDetails.ShipmentRateDetail.GrandTotal.Amount

			services[p.ServiceType] = fmt.Sprintf("%.2f", a)

		}
		rates[p.Id] = map[string]map[string]string{"FEDEX": services}
	}

	return rates, nil
}
