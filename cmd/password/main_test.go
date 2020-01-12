package main

import (
	"testing"
)

const id = "0fc6742e-c35b-47a1-a80b-0037bdcd4e8e"
const val = "Pass123!"

func TestHandleOk(t *testing.T) {
	if err := Handle(Password{Id: id, Password: val}); err != nil {
		t.Fail()
	}
}

func TestHandleBadId(t *testing.T) {
	if err := Handle(Password{Id: "", Password: val}); err == nil {
		t.Fail()
	}
}

func TestHandleBadVal(t *testing.T) {
	if err := Handle(Password{Id: id, Password: "Pass1234!"}); err == nil {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
