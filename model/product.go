package model

type Product struct {
	ID            string   `json:"id";repo:"id,hash"`
	SKU           string   `json:"sku";repo:"sku,sort"`
	Category      string   `json:"category";repo:"category"`
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
