package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Product struct {
	Id            string   `json:"id"`
	Session       string   `json:"session"`
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
	ShipFrom      string   `json:"ship_from"`
}

func (p *Product) Unmarshal(s string) error {
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return err
	} else if len(p.Name) < 3 {
		return fmt.Errorf("bad name [%s], must be at least 3 characters in length", p.Name)
	} else if p.PriceInteger == "" {
		return fmt.Errorf("bad price (integer) [%s]", p.PriceInteger)
	} else if p.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		p.Id = id.String()
		return nil
	}
}
