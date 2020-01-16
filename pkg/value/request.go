package value

type Request struct {
	Case       string      `json:"case"`
	Table      string      `json:"table"`
	Handler    string      `json:"handler"`
	Ids        []string    `json:"ids"`
	Result     interface{} `json:"result,omitempty"`
	BrandId    string      `json:"brand_id,omitempty"`
	Expression string      `json:"expression,omitempty"`
}
