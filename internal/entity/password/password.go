package password

import (
	"golang.org/x/crypto/bcrypt"
	"os"
)

// Data structure for persisting and retrieving a users (encrypted) password. The user entity maintains a 1-1 FK
// relationship with the UserPassword entity for referential integrity, where user.PasswordId == userPassword.Id.
type Entity struct {
	ID      string `json:"id"`
	Encoded string `json:"password"`
	Decoded []byte `json:"password_dec"`
}

var handler = "hcPasswordHandler"
var table = os.Getenv("USER_PASSWORD_TABLE")

func (e *Entity) Id() *string {
	return &e.ID
}

func (e *Entity) Validate() error {
	return bcrypt.CompareHashAndPassword(e.Payload(), e.Decoded)
}

func (e *Entity) Payload() []byte {
	return []byte(e.Encoded)
}

func (*Entity) Handler() *string {
	return &handler
}

func (*Entity) Name() *string {
	return &table
}
