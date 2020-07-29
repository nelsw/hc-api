package product

type Entity struct {
	Id       string   `json:"id"`
	Owner    string   `json:"owner"`
	Brand    string   `json:"brand"`
	Category string   `json:"category"`
	Name     string   `json:"name"`
	Summary  string   `json:"summary"`
	Image    string   `json:"image"`
	Options  []Option `json:"options"`
}

type Option struct {
	Id      string   `json:"id"`
	Parent  string   `json:"parent"`
	Price   int64    `json:"price"`   // 7900 = $79.00, stripe thinks it makes cents
	Weight  int      `json:"weight"`  // 170 = 1.7, to avoid decimals entirely
	Label   string   `json:"label"`   // oz, lb, kilo, ton, w/e
	Stock   int      `json:"stock"`   // quantity available
	Address string   `json:"address"` // shipping departure location
	Images  []string `json:"images"`  // urls
}

func (e *Entity) Validate() error {
	return nil
}
