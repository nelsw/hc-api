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
	"github.com/dgrijalva/jwt-go"
	"os"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/profile"
)

var table = os.Getenv("TABLE")

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var e profile.Entity

	_, err := apigwp.Request(r, &e)
	if err != nil {
		return apigwp.Response(400, err)
	}

	token, ok := r.QueryStringParameters["token"]
	if !ok {
		return apigwp.Response(400, fmt.Errorf("no token provided"))
	}

	switch r.Path {

	case "save":
		inspection := events.APIGatewayProxyRequest{
			Path:                  "inspect",
			QueryStringParameters: map[string]string{"token": token},
		}
		code, body := client.CallIt(inspection, "tokenHandler")
		if code != 200 {
			return apigwp.Response(code, body)
		}

		response := events.APIGatewayProxyResponse{}
		_ = json.Unmarshal([]byte(body), &response)

		if response.StatusCode != 200 {
			return apigwp.Response(response.StatusCode, response.Body)
		}

		claims := jwt.StandardClaims{}
		if err := json.Unmarshal([]byte(response.Body), &claims); err != nil {
			return apigwp.Response(400, err)
		}

		m := map[string]interface{}{
			"id":      claims.Id,
			"table":   table,
			"type":    "*profile.Entity",
			"keyword": "save",
			"result":  &e,
		}
		return apigwp.Response(client.CallIt(&m, "repoHandler"))
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
