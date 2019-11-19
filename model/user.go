// The structures and functions made available in this file define the user data model and foreign entity relationships.
package model

import (
	"fmt"
	"regexp"
	"unicode"
)

// No email regex is perfect, but this one is close.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Primary user object for the domain, visible to client and server. With the exception of the Email field, which
// represents its plain text value, each property is a reference to a unique ID, or collection of unique ID's.
type User struct {
	Email      string   `json:"email"`
	PasswordId string   `json:"password_id"`
	ProfileId  string   `json:"profile_id"`
	AddressIds []string `json:"address_ids"`
	ProductIds []string `json:"product_ids"`
	OrderIds   []string `json:"order_ids"`
	SaleIds    []string `json:"sale_ids"`
	Session    string   `json:"session"`
}

// Used for login and registration use cases.
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Data structure for persisting and retrieving a users (encrypted) password. The user entity maintains a 1-1 FK
// relationship with the UserPassword entity for referential integrity, where user.PasswordId == userPassword.Id.
type UserPassword struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

// UserProfile also promotes separation of concerns by decoupling user profile details from the primary user entity. IF
// UserProfile.EmailOld != UserProfile.EmailNew, AND User.Email == UserProfile.EmailOld, THEN we must prompt the user to
// confirm new email address. IF UserProfile.Password1 is not blank AND UserProfile.Password2 is not blank AND valid AND
// UserProfile.Password1 == UserProfile.Password2, then we update the UserPassword entity and return OK.
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

// Used to update an existing user item in an Amazon DynamoDB table.
// SET - modify or add item attributes
// REMOVE - delete attributes from an item
// ADD - update numbers and sets
// DELETE - remove elements from a set
type UserUpdate struct {
	Val        string `json:"val"`
	Expression string `json:"expression"`
	Session    string `json:"session"`
}

// Validates the UserCredentials entity by confirming that both the email and password values are valid.
// Allows email addresses with third party domains and any extension.
func (uc *UserCredentials) Validate() error {
	if emailRegex.MatchString(uc.Email) == false {
		return fmt.Errorf("bad email [%s]", uc.Email)
	} else if err := IsPasswordValid(uc.Password); err != nil {
		return err
	} else {
		return nil
	}
}

// Validates the UserPassword entity by confirming that both the password and id values are valid.
func (up *UserPassword) Validate() error {
	return IsPasswordValid(up.Password)
}

// The following is an adaptation of https://stackoverflow.com/a/25840157
func IsPasswordValid(s string) error {
	var number, upper, special bool
	length := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			length++
		case unicode.IsUpper(c):
			upper = true
			length++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			length++
		case unicode.IsLetter(c) || c == ' ':
			length++
		default:
			// do not increment length for unrecognized characters
		}
	}
	if length < 8 || length > 24 {
		return fmt.Errorf("bad password, must contain 8-24 characters")
	} else if number == false {
		return fmt.Errorf("bad password, must contain at least 1 number")
	} else if upper == false {
		return fmt.Errorf("bad password, must contain at least 1 uppercase letter")
	} else if special == false {
		return fmt.Errorf("bad password, must contain at least 1 special character")
	} else {
		return nil
	}
}
