package product

import "os"

type Aggregate struct {
	Id          string   `json:"id"`
	Sku         string   `json:"sku"`
	AddressId   string   `json:"address_id"`
	OwnerId     string   `json:"owner_id"`
	Img         string   `json:"img"`
	Category    string   `json:"category"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	ImageSet    []string `json:"image_set"`
	Quantity    string   `json:"quantity"`
	Stock       string   `json:"stock"`
	Deleted     string   `json:"deleted,omitempty"`
	// packaging details (calc shipping rates)
	Unit    string  `json:"unit"` // LB
	Weight  float32 `json:"weight"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Length  int     `json:"length"`
	Session string  `json:"session"`
}

var table = os.Getenv("PRODUCT_TABLE")
