package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Request struct {
	Command      string `json:"command"`
	UserId       string `json:"user_id"`
	SourceIp     string `json:"source_ip"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_KEY"))
var regex = regexp.MustCompile(`(.*token=)(.*)(;.*)`)

func keyFunc(token *jwt.Token) (interface{}, error) { return jwtKey, nil }

func newToken(request Request, duration time.Duration, name string) (string, error) {
	now := time.Now()
	expiry := now.Add(duration)
	if token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  request.SourceIp,
			ExpiresAt: expiry.Unix(),
			Id:        request.UserId,
			IssuedAt:  now.Unix(),
		},
	}).SignedString(jwtKey); err != nil {
		return "", nil
	} else {
		cookie := &http.Cookie{Name: name + "-token", Value: token, Expires: expiry, HttpOnly: false}
		return cookie.String(), nil
	}
}

func validate(token, ip, name string) (string, error) {
	claims := &Claims{}
	if tkn, err := jwt.ParseWithClaims(regex.ReplaceAllString(token, `$2`), claims, keyFunc); err != nil {
		return "", err // Either the token expired or the signature doesn't match.
	} else if !tkn.Valid {
		return "", fmt.Errorf("bad [%s] token=[%v] claims=[%v]", name, tkn, claims)
	} else if claims.Audience != ip {
		return "", fmt.Errorf("bad [%s] ips got=[%s] want=[%s] claims=[%v]", name, ip, claims.Audience, claims)
	} else {
		return claims.Id, nil
	}
}

func Handle(request Request) (string, error) {
	fmt.Printf("REQUEST=[%v]\n", request)
	switch request.Command {
	case "validate":
		if _, err := validate(request.AccessToken, request.SourceIp, "access"); err != nil {
			return "", err
		} else {
			return validate(request.RefreshToken, request.SourceIp, "refresh")
		}
	case "refresh":
		return newToken(request, 30*time.Minute, "refresh")
	case "access":
		return newToken(request, 24*time.Hour, "access")
	}
	return "", fmt.Errorf("bad request")
}

func main() {
	lambda.Start(Handle)
}
