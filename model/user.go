package model

// Primary user object for the domain, visible to client and server.
// With the exception of the Email field, which represents its plain text value,
// each property is a reference to a unique ID, or collection of unique ID's.
type User struct {
	Email      string   `json:"email"`
	PasswordId string   `json:"-"`
	ProfileId  string   `json:"profile_id"`
	AddressIds []string `json:"address_ids"`
	ProductIds []string `json:"product_ids"`
	OrderIds   []string `json:"order_ids"`
	SaleIds    []string `json:"sale_ids"`
}

func (user *User) Validate() error {
	if err := IsEmailValid(user.Email); err != nil {
		return err
	} else if err := IsIdValid(user.PasswordId); err != nil {
		return err
	} else if err := IsIdValid(user.ProfileId); err != nil {
		return err
	} else {
		return nil
	}
}
