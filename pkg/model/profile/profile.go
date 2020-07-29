package profile

type Entity struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Img       string `json:"img"`
}

func (*Entity) Validate() error {
	return nil
}
