package main

import (
	"testing"
)

func TestHandle(t *testing.T) {
	if str, err := Handle(Request{
		Command:  "access",
		UserId:   "0123456789",
		SourceIp: "127.0.0.1",
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(str)
	}
}
