package class

type Request struct {
	Expression string `json:"expression,omitempty"`

	Command string   `json:"command"`
	Table   string   `json:"table"`
	Ids     []string `json:"ids"`

	BrandId string `json:"brand_id,omitempty"`

	Data    interface{}   `json:"data,omitempty"`
	DataArr []interface{} `json:"data_arr,omitempty"`
}
