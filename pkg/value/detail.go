package value

type Detail struct {
	Sku         string `json:"sku"`
	Img         string `json:"img"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Other       string `json:"other"`
}
