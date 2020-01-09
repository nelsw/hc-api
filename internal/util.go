package internal

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func NewToken(id, ip string) string {
	type Claims struct {
		Id string `json:"id"`
		Ip string `json:"ip"`
		jwt.StandardClaims
	}
	expiry := time.Now().Add(30 * time.Minute)
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{Id: id, Ip: ip, StandardClaims: jwt.StandardClaims{ExpiresAt: expiry.Unix()}})
	str, _ := tkn.SignedString(jwtKey)
	cookie := &http.Cookie{Name: "token", Value: str, Expires: expiry, HttpOnly: false}
	return cookie.String()
}
