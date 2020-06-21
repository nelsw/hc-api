package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"regexp"
	"sam-app/pkg/factory/apigwp"
	"time"
)

var regex = regexp.MustCompile(`(.*)(token=)(.*)(;.*)`)
var jwtKey = []byte(os.Getenv("JWT_KEY"))

func keyFunc(_ *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}

func authenticate(token string, claims *jwt.StandardClaims) error {
	if jwtToken, err := jwt.ParseWithClaims(token, claims, keyFunc); err != nil || !jwtToken.Valid {
		return fmt.Errorf("bad token, invalid segments or expired\n")
	} else {
		return nil
	}
}

func issue(claims *jwt.StandardClaims) string {
	claims.IssuedAt = time.Now().Unix()
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = time.Unix(claims.IssuedAt, 0).Add(time.Hour * 24).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := token.SignedString(jwtKey)
	cookie := &http.Cookie{
		Name:    "token",
		Value:   str,
		Expires: time.Unix(claims.ExpiresAt, 0),
	}

	return cookie.String()
}

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	apigwp.LogRequest(r)

	switch r.Path {

	case "authenticate":
		if token, ok := r.Headers["Authorize"]; ok {
			tokenString := regex.ReplaceAllString(token, `$3`)
			claims := jwt.StandardClaims{}
			if err := authenticate(tokenString, &claims); err != nil {
				return apigwp.Response(401, err)
			}
			newToken := issue(&claims)
			return apigwp.ProxyResponse(200, map[string]string{"Authorize": token}, newToken)
		}

	case "authorize":
		claims := jwt.StandardClaims{}
		if err := json.Unmarshal([]byte(r.Body), &claims); err != nil {
			return apigwp.Response(400, err)
		}
		token := issue(&claims)
		return apigwp.ProxyResponse(200, map[string]string{"Authorize": token}, &token)

	case "inspect":
		if token, ok := r.Headers["Authorize"]; ok {
			tokenString := regex.ReplaceAllString(token, `$3`)
			claims := jwt.StandardClaims{}
			if err := authenticate(tokenString, &claims); err != nil {
				return apigwp.Response(401, err)
			}
			return apigwp.ProxyResponse(200, map[string]string{"Authorize": issue(&claims)}, &claims)
		}
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
