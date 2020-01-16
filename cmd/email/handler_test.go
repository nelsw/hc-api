package email

import (
	"encoding/json"
	"hc-api/pkg/value"
	"testing"
)

func TestSendEmail(t *testing.T) {

	b, _ := json.Marshal(&value.Email{
		To:       "connorvanelswyk@gmail.com",
		Subject:  "test email",
		Body:     "",
		Code:     "wat",
		Template: "password-reset.html",
	})

	if err := SendEmail(string(b)); err != nil {
		t.Error(err)
	}

}
