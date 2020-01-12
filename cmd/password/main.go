package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
	. "hc-api/service"
	"os"
)

// Data structure for persisting and retrieving a users (encrypted) password. The user entity maintains a 1-1 FK
// relationship with the UserPassword entity for referential integrity, where user.PasswordId == userPassword.Id.
type Password struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

var table = os.Getenv("USER_PASSWORD_TABLE")

func Handle(p Password) error {
	got := []byte(p.Password)
	_ = FindOne(&table, &p.Id, &p)
	want := []byte(p.Password)
	return bcrypt.CompareHashAndPassword(want, got)
}

func main() {
	lambda.Start(Handle)
}
