package main

import (
	"testing"
)

const email = "connor@wiesow.com"
const userId = "638b13ef-ab84-410a-abb0-c9fd5da45c62"
const passwordVal = "Pass123!"

func TestHandleOk(t *testing.T) {
	if _, err := Handle(ClientCredential{Email: email, Password: passwordVal}); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestHandleBadEmail(t *testing.T) {
	if _, err := Handle(ClientCredential{Email: "", Password: passwordVal}); err == nil {
		t.Fail()
	}
}

func TestHandleBadPassword(t *testing.T) {
	if _, err := Handle(ClientCredential{Email: email, Password: "Pass1234!"}); err == nil {
		t.Fail()
	}
}

func TestHandleBadLength(t *testing.T) {
	if _, err := Handle(ClientCredential{Email: email, Password: "Pass12!"}); err != LengthErr {
		t.Fail()
	}
}

func TestHandleBadNumbers(t *testing.T) {
	if _, err := Handle(ClientCredential{Email: email, Password: "Pass!!!!"}); err != NumberErr {
		t.Fail()
	}
}

func TestHandleBadSpecialChar(t *testing.T) {
	if _, err := Handle(ClientCredential{Email: email, Password: "Pass1234"}); err != CharErr {
		t.Fail()
	}
}

func TestHandleBadUppercase(t *testing.T) {
	if _, err := Handle(ClientCredential{Email: email, Password: "pass123!"}); err != CaseErr {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
