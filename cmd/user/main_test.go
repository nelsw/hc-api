package main

import (
	"sam-app/pkg/model/user"
	"sam-app/test"
	"testing"
)

func TestHandleAdd(t *testing.T) {
	r := user.Request{
		"add",
		test.UserId,
		"address_ids",
		[]string{test.AddressId},
	}
	if err := Handle(r); err != nil {
		t.Fatal(err)
	}
}

func TestHandleDelete(t *testing.T) {
	r := user.Request{
		"delete",
		test.UserId,
		"address_ids",
		[]string{test.AddressId},
	}
	if err := Handle(r); err != nil {
		t.Fatal(err)
	}
}

func TestHandleBadOp(t *testing.T) {
	if err := Handle(user.Request{}); err != ErrBadOp {
		t.Fatal(err)
	}
}

// for code coverage purposes only
func TestHandleRequest(t *testing.T) {
	go main()
}
