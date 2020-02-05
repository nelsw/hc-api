package request

type Entity struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Table      string            `json:"table"`
	Ids        []string          `json:"ids,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Keyword    string            `json:"keyword,omitempty"`
	Result     interface{}       `json:"result,omitempty"`
	Results    []interface{}     `json:"results,omitempty"`
}
