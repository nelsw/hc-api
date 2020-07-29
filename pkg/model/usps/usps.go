package usps

import (
	"encoding/xml"
)

type Address struct {
	Id     string `json:"id" xml:"ID,attr"`
	Unit   string `json:"unit" xml:"Address1"`
	Street string `json:"street" xml:"Address2"`
	City   string `json:"city" xml:"City"`
	State  string `json:"state" xml:"State"`
	Zip5   string `json:"zip_5" xml:"Zip5"`
	Zip4   string `json:"zip_4" xml:"Zip4"`
}

type Request struct {
	Address  Address          `json:"address"`
	Packages []PackageRequest `json:"packages"`
}

type AddressValidateRequest struct {
	UserId   string  `xml:"USERID,attr"`
	Revision string  `xml:"Revision"`
	Address  Address `xml:"Address"`
}

type AddressValidateResponse struct {
	Address Address `xml:"Address"`
}

type RateV4Request struct {
	UserId   string           `xml:"USERID,attr"`
	Revision string           `xml:"Revision"`
	Packages []PackageRequest `xml:"Package"`
}

type RateV4Response struct {
	Packages []PackageResponse `xml:"Package"`
}

type PackageRequest struct {
	XMLName    xml.Name `xml:"Package"`
	Id         string   `xml:"ID,attr" json:"id"` // product id
	Service    string   `xml:"Service" json:"-"`
	ZipFrom    string   `xml:"ZipOrigination" json:"zip_from"`
	ZipTo      string   `xml:"ZipDestination" json:"zip_to"`
	Pounds     int      `xml:"Pounds" json:"pounds"`
	Ounces     float32  `xml:"Ounces" json:"ounces"`
	Container  string   `xml:"Container" json:"-"`
	Machinable string   `xml:"Machinable" json:"-"`
}

type PackageResponse struct {
	PackageRequest
	Postage Postage `xml:"Postage" json:"postage"`
}

type Postage struct {
	Id    string `xml:"CLASSID,attr" json:"id"`
	Type  string `xml:"MailService" json:"type"`
	Price string `xml:"Rate" json:"price"`
}
