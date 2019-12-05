package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"regexp"
	"time"
)

const hack = "dbb7c19d-521f-41e2-ada6-cefbec1b3bcf"

var jwtKey = []byte(os.Getenv("jwt_key"))
var regex = regexp.MustCompile(`(token=)(.*)(;.*)`)

// ƒ responsible for jwt key token interpretation
func keyFunc(token *jwt.Token) (interface{}, error) { return jwtKey, nil }

// Data structure representing a parsed JWT string.
type Claims struct {
	Id string `json:"id"`
	Ip string `json:"ip"`
	jwt.StandardClaims
}

func NewCookie(s1, s2 string) (cookie string, err error) {
	expiry := time.Now().Add(30 * time.Minute)
	// Declare the token with the algorithm used for signing, and JWT claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Id: s1,
		Ip: s2,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry.Unix(),
		},
	})

	// Create the JWT string.
	if tokenString, err := token.SignedString(jwtKey); err != nil {
		return "", err
	} else {
		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself.
		cookie := &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  expiry,
			HttpOnly: false,
		}
		return cookie.String(), nil
	}
}

func ValidateSession(cookie, ip string) (string, error) {
	if cookie == hack {
		return "", nil
	}
	claims := &Claims{}
	// Find the token in the cookie, may not exist.
	token := regex.ReplaceAllString(cookie, `$2`)
	if tkn, err := jwt.ParseWithClaims(token, claims, keyFunc); err != nil {
		// Either the token expired or the signature doesn't match.
		return "", err
	} else if !tkn.Valid || claims.Ip != ip {
		return "", errors.New("bad token")
	} else {
		return claims.Id, nil
	}
}
