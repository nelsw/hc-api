package util

import (
	"sam-app/pkg/model/product"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	if err := ValidateEmail("foo"); err == nil {
		t.Fail()
	}
}

func TestValidateEmailErr(t *testing.T) {
	if err := ValidateEmail("foo@gmail.com"); err != nil {
		t.Fail()
	}
}

func TestValidatePassword(t *testing.T) {
	if err := ValidatePassword("Pass123!"); err != nil {
		t.Fail()
	}
}

func TestValidatePasswordLengthErr(t *testing.T) {
	if err := ValidatePassword("Pass123"); err == nil {
		t.Fail()
	}
}

func TestValidatePasswordNumberErr(t *testing.T) {
	if err := ValidatePassword("Pass!!!!"); err == nil {
		t.Fail()
	}
}

func TestValidatePasswordCaseErr(t *testing.T) {
	if err := ValidatePassword("pass123!"); err == nil {
		t.Fail()
	}
}

func TestValidatePasswordCharErr(t *testing.T) {
	if err := ValidatePassword("Pass1234"); err == nil {
		t.Fail()
	}
}

func TestTypeName(t *testing.T) {
	s := TypeName(product.Entity{})
	t.Log(s) // first step towards generics
}
