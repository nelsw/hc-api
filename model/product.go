package model

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
