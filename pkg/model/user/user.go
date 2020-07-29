package user

type Entity struct {
	Id       string   `json:"id"`
	Products []string `json:"products"`
	Orders   []string `json:"orders"`
}
