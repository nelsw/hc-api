// UserProfile is exactly what it appears to be, a user item domain model validation.
// It also promotes separation of concerns by decoupling user item details from the primary user validation. IF
// UserProfile.EmailOld != UserProfile.EmailNew, AND User.Credential == UserProfile.EmailOld, THEN we must prompt the user to
// confirm new email validation. IF UserProfile.Password1 is not blank AND UserProfile.Password2 is not blank AND valid AND
// UserProfile.Password1 == UserProfile.Password2, then we update the UserPassword validation and return OK.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"hc-api/pkg/client/faas/client"
	"hc-api/pkg/client/repo"
	"hc-api/pkg/factory/apigwp"
	"hc-api/pkg/model/profile"
	"hc-api/pkg/model/token"
)

var InvalidRequest = fmt.Sprintf("bad request\n")

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var e profile.Proxy

	ip, err := apigwp.Request(r, &e)
	if err != nil {
		return apigwp.Response(400, err)
	}

	op := e.Subject

	e.Subject = "authenticate"
	e.SourceIp = ip // we set the ip here for CORS prevention

	tkn := token.Entity{e.Value, token.Error{}}
	out, err := client.Invoke(&tkn)
	if err != nil {
		return apigwp.Response(500, err)
	}

	_ = json.Unmarshal(out, &tkn)
	if tkn.Error.Msg != "" {
		return apigwp.Response(402, &tkn)
	}

	switch op {

	case "save":
		newProfile := e.Id == ""
		if newProfile {
			s, _ := uuid.NewUUID()
			e.Id = s.String()
		}
		if err := repo.SaveOne(&e.Entity); err != nil {
			return apigwp.Response(500, err)
		}
		if newProfile {
			// todo
		}
		return apigwp.Response(200, &e)
	}

	return apigwp.Response(400, InvalidRequest)
}

func main() {
	lambda.Start(Handle)
}
