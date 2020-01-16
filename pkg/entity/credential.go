package entity

import (
	"hc-api/pkg/value"
	"os"
	"strings"
)

// Used for login and registration use cases.
// Result structure for securely associating user entities with their credential and password.
type Credential struct {
	Id         string `json:"id"` // email address
	UserId     string `json:"user_id"`
	PasswordId string `json:"password_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

var credentialHandler = "hcCredentialHandler"
var credentialTable = os.Getenv("CREDENTIAL_TABLE")

func (e *Credential) Ids() []string {
	return []string{strings.ToLower(e.Email)}
}

func (e *Credential) Payload() []byte {
	return []byte(e.Password)
}

func (*Credential) Function() *string {
	return &credentialHandler
}

func (*Credential) Table() *string {
	return &credentialTable
}

// Confirm we have either a valid username or email for an ID and the password meets modern strength criteria.
func (e *Credential) Validate() error {
	if err := value.ValidateEmail(e.Email); err != nil {
		return err
	} else if err := value.ValidatePassword(e.Password); err != nil {
		return err
	} else {
		return nil
	}
}
