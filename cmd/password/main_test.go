package main

import (
	"testing"
)

const (
	id  = "0fc6742e-c35b-47a1-a80b-0037bdcd4e8e"
	val = "Pass123!"
)

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
	if err := Handle(Password{Id: id, Password: ""}); err == nil {
		t.Fail()
	}
}

func TestHandleBadLength(t *testing.T) {
	if err := Handle(Password{Id: id, Password: "Pass12!"}); err != lengthError {
		t.Fail()
	}
}

func TestHandleBadNumbers(t *testing.T) {
	if err := Handle(Password{Id: id, Password: "Pass!!!!"}); err != numberError {
		t.Fail()
	}
}

func TestHandleBadSpecialChar(t *testing.T) {
	if err := Handle(Password{Id: id, Password: "Pass1234"}); err != specialError {
		t.Fail()
	}
}

func TestHandleBadUppercase(t *testing.T) {
	if err := Handle(Password{Id: id, Password: "pass123!"}); err != upperError {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
