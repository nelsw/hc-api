package apigw

import (
	"encoding/json"
	"fmt"
	"testing"
)

type ExampleStruct struct {
	Id, Ip string
}

func TestOk(t *testing.T) {
	if r, _ := Ok(); r.StatusCode != 200 || r.Body != "" {
		t.Fail()

	}
}

func TestOkString(t *testing.T) {
	if r, _ := Ok("Hello World!"); r.StatusCode != 200 || r.Body != "Hello World!" {
		t.Fail()
	}
}

func TestOkStruct(t *testing.T) {
	want := ExampleStruct{Id: "123456790", Ip: "127.0.0.1"}
	var got ExampleStruct
	r, _ := Ok(want)
	_ = json.Unmarshal([]byte(r.Body), &got)
	if got != want || r.StatusCode != 200 {
		t.Fail()
	}
}

func TestBadAuth(t *testing.T) {
	if r, _ := BadAuth(fmt.Errorf("error msg")); r.StatusCode != 401 || r.Body != "error msg" {
		t.Fail()
	}
}

func TestBadRequest(t *testing.T) {
	if r, _ := BadRequest(); r.StatusCode != 400 || r.Body != "" {
		t.Fail()
	}
}

func TestBadRequestString(t *testing.T) {
	if r, _ := BadRequest("error msg"); r.StatusCode != 400 || r.Body != "error msg" {
		t.Fail()
	}
}

func TestBadRequestError(t *testing.T) {
	if r, _ := BadRequest(fmt.Errorf("error msg")); r.StatusCode != 400 || r.Body != "error msg" {
		t.Fail()
	}
}
