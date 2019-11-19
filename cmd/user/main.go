package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	response "github.com/nelsw/hc-util/aws"
	"golang.org/x/crypto/bcrypt"
	"hc-api/model"
	"hc-api/repo"
	"hc-api/service"
	"log"
	"net/http"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("request: [%v]", r)

	cmd := r.QueryStringParameters["cmd"]

	switch cmd {

	case "register":
		// todo - create User, UserPassword, and UserProfile entities ... also verify email address.
		return response.New().Code(http.StatusNotImplemented).Build()

	case "login":
		var uc model.UserCredentials
		if err := json.Unmarshal([]byte(r.Body), &uc); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := uc.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if user, err := repo.FindUserByEmail(&uc.Email); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if up, err := repo.FindUserPasswordById(&user.PasswordId); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := bcrypt.CompareHashAndPassword([]byte(up.Password), []byte(uc.Password)); err != nil {
			return response.New().Code(http.StatusUnauthorized).Build()
		} else if cookie, err := service.NewCookie(user.Email); err != nil {
			return response.New().Code(http.StatusInternalServerError).Build()
		} else {
			user.Session = cookie
			user.OrderIds = nil
			user.SaleIds = nil
			return response.New().Code(http.StatusOK).Data(&user).Build()
		}

	case "update":
		var uu model.UserUpdate
		if err := json.Unmarshal([]byte(r.Body), &uu); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if email, err := service.Validate(uu.Session); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := repo.UpdateUser(&email, &uu.Expression, &uu.Val); err != nil {
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
