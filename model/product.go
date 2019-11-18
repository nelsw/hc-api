package model

import "fmt"

type Product struct {
	Id            string   `json:"id"`
	Sku           string   `json:"sku"`
	Category      string   `json:"category"`
	Subcategory   string   `json:"subcategory"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PriceInteger  string   `json:"price_integer"`
	PriceFraction string   `json:"price_fraction"`
	Quantity      string   `json:"quantity"`
	Unit          string   `json:"unit"`
	Owner         string   `json:"owner"`
	ImageSet      []string `json:"image_set"`
}

func (p *Product) Validate() error {
	length := len(p.Name)
	if length < 3 {
		return fmt.Errorf("bad name [%s], must be at least 3 characters in length", p.Name)
	} else if p.PriceInteger == "" {
		return fmt.Errorf("bad price (integer) [%s]", p.PriceInteger)
	} else {
		return nil
	}
}
