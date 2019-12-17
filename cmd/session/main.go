package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	response "github.com/nelsw/hc-util/aws"
	"net/http"
	"os"
	"regexp"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))
var regex = regexp.MustCompile(`(token=)(.*)(;.*)`)

// ƒ responsible for jwt key token interpretation
func keyFunc(token *jwt.Token) (interface{}, error) { return jwtKey, nil }

// Data structure representing a parsed JWT string.
type Claims struct {
	Id string `json:"id"`
	Ip string `json:"ip"`
	jwt.StandardClaims
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	ip := r.QueryStringParameters["ip"]
	fmt.Printf("REQUEST [%s]: [%v]", cmd, r)

	switch cmd {

	case "validate":
		claims := &Claims{}
		session := r.QueryStringParameters["session"]
		token := regex.ReplaceAllString(session, `$2`)
		if tkn, err := jwt.ParseWithClaims(token, claims, keyFunc); err != nil {
			// Either the token expired or the signature doesn't match.
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if !tkn.Valid || claims.Ip != ip {
			return response.New().Code(http.StatusUnauthorized).Text("bad token").Build()
		} else {
			return response.New().Code(http.StatusOK).Text(claims.Id).Build()
		}

	case "create":
		id := r.QueryStringParameters["id"]
		expiry := time.Now().Add(30 * time.Minute)
		// Declare the token with the algorithm used for signing, and JWT claims.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
			Id: id,
			Ip: ip,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiry.Unix(),
			},
		})
		// Create the JWT string.
		if tokenString, err := token.SignedString(jwtKey); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			// Finally, we set the client cookie for "token" as the JWT we just generated
			// we also set an expiry time which is the same as the token itself.
			cookie := &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				Expires:  expiry,
				HttpOnly: false,
			}
			return response.New().Code(http.StatusOK).Text(cookie.String()).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
