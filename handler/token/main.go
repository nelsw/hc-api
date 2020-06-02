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
var jwtKey = os.Getenv("JWT_KEY")

func keyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(jwtKey), nil
}

func issue(claims *jwt.StandardClaims) string {
	claims.IssuedAt = time.Now().Unix()
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = time.Unix(claims.IssuedAt, 0).Add(time.Second * 24).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := token.SignedString([]byte(jwtKey))
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
		if token, ok := r.QueryStringParameters["token"]; ok {
			tokenString := regex.ReplaceAllString(token, `$3`)
			claims := jwt.StandardClaims{}
			if jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc); err != nil {
				return apigwp.Response(401, err) // Either the token expired or signatures do not match.
			} else if !jwtToken.Valid {
				return apigwp.Response(401, fmt.Errorf("bad token, invalid segments or expired\n"))
			} else {
				return apigwp.Response(200, issue(&claims))
			}
		}

	case "authorize":

		claims := jwt.StandardClaims{}

		if err := json.Unmarshal([]byte(r.Body), &claims); err != nil {
			return apigwp.Response(400, err)
		}

		return apigwp.Response(200, issue(&claims))
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
