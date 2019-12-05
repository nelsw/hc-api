package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	response "github.com/nelsw/hc-util/aws"
	"golang.org/x/crypto/bcrypt"
	"hc-api/model"
	"hc-api/service"
	"net/http"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	fmt.Printf("REQUEST [%s]: [%v]", cmd, r)

	switch cmd {

	case "register":
		// todo - create User, UserPassword, and UserProfile entities ... also verify email address.
		return response.New().Code(http.StatusNotImplemented).Build()

	case "login":
		var uc model.UserCredentials
		if err := uc.Unmarshal(r.Body); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if ue, err := service.FindUserEmailById(&uc.Email); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if u, err := service.FindUserById(&ue.UserId); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if up, err := service.FindUserPasswordById(&u.PasswordId); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if err := bcrypt.CompareHashAndPassword([]byte(up.Password), []byte(uc.Password)); err != nil {
			return response.New().Code(http.StatusUnauthorized).Build()
		} else if cookie, err := service.NewCookie(u.Id, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusInternalServerError).Build()
		} else {
			u.Session = cookie
			u.OrderIds = nil
			u.SaleIds = nil
			u.Id = ""
			u.PasswordId = ""
			return response.New().Code(http.StatusOK).Data(&u).Build()
		}

	case "update":
		var uu model.UserUpdate
		if err := json.Unmarshal([]byte(r.Body), &uu); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if id, err := service.ValidateSession(uu.Session, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := service.UpdateUser(&id, &uu.Expression, &uu.Val); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Text(string([]byte(`{"success":""}`))).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
