package credential

import (
	"hc-api/internal/util/validate"
	"os"
	"strings"
)

type Aggregate struct {
	Value
	Entity
}

// Used for login and registration use cases.
type Value struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Data structure for securely associating user entities with their credential and password.
type Entity struct {
	Id         string `json:"id"` // email address
	UserId     string `json:"user_id"`
	PasswordId string `json:"password_id"`
}

var handler = "hcCredentialHandler"
var table = os.Getenv("CREDENTIAL_TABLE")

func (e *Aggregate) Id() *string {
	s := strings.ToLower(e.Email)
	return &s
}

func (e *Aggregate) Payload() []byte {
	return []byte(e.Password)
}

func (*Aggregate) Handler() *string {
	return &handler
}

func (*Aggregate) Name() *string {
	return &table
}

// Confirm we have either a valid username or email for an ID and the password meets modern strength criteria.
func (e *Aggregate) Validate() error {
	if err := validate.Email(e.Email); err != nil {
		return err
	} else if err := validate.Password(e.Password); err != nil {
		return err
	} else {
		return nil
	}
}
