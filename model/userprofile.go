package model

type UserProfile struct {
	Id        string `json:"id"`
	EmailOld  string `json:"email_old"`
	EmailNew  string `json:"email_new"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Password1 string `json:"password_1"`
	Password2 string `json:"password_2"`
}
