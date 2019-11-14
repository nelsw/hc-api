package dynamo

import (
	"hc-api/model"
	"hc-api/repo"
	"testing"
)

func TestGetUser(t *testing.T) {
	user := model.User{
		Email:    "connor@wiesow.com",
		Password: "Pass123!",
	}
	err := repo.FindUser(&user)
	if err != nil {
		t.Errorf("could not get test user. details: %v", err)
	}
}
