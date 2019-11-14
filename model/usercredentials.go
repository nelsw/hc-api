package model

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uc *UserCredentials) Validate() error {
	if err := IsEmailValid(uc.Email); err != nil {
		return err
	} else if err := IsPasswordValid(uc.Password); err != nil {
		return err
	} else {
		return nil
	}
}
