package class

type Address struct {
	Street string `json:"street"`
	Unit   string `json:"unit,omitempty"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip5   string `json:"zip_5"`
	Zip4   string `json:"zip_4,omitempty"`
}
