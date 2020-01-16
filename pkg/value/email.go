package value

type Email struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Code     string `json:"code"`
	Template string `json:"template"`
}
