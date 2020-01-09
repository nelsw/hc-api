// UserProfile is exactly what it appears to be, a user profile domain model entity.
// It also promotes separation of concerns by decoupling user profile details from the primary user entity. IF
// UserProfile.EmailOld != UserProfile.EmailNew, AND User.Email == UserProfile.EmailOld, THEN we must prompt the user to
// confirm new email address. IF UserProfile.Password1 is not blank AND UserProfile.Password2 is not blank AND valid AND
// UserProfile.Password1 == UserProfile.Password2, then we update the UserPassword entity and return OK.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/apigw"
	. "hc-api/service"
	"os"
)

var table = os.Getenv("USER_PROFILE_TABLE")

type Profile struct {
	Id        string   `json:"id"`
	BrandIds  []string `json:"brand_ids"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Phone     string   `json:"phone"`
}

type ProfileRequest struct {
	Command string  `json:"command"`
	Session string  `json:"session"`
	Profile Profile `json:"profile"`
	Id      string  `json:"id"`
}

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var r ProfileRequest
	body := request.Body
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return apigw.BadRequest(err)
	}

	fmt.Printf("REQUEST   [%s]\n", body)

	ip := request.RequestContext.Identity.SourceIP
	if _, err := Invoke().Handler("Session").Session(r.Session).IP(ip).CMD("validate").Post(); err != nil {
		return apigw.BadAuth(err)
	}

	switch r.Command {

	case "save":
		if err := Put(&r.Profile, &table); err != nil {
			return apigw.BadRequest(err)
		} else {
			return apigw.Ok()
		}

	case "find":
		if err := FindOne(&table, &r.Id, &r.Profile); err != nil {
			return apigw.BadRequest(err)
		} else {
			return apigw.Ok(&r.Profile)
		}

	default:
		return apigw.BadRequest()
	}
}

func main() {
	lambda.Start(Handle)
}
