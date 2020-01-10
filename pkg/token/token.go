package token

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

type Claims struct {
	Id string `json:"id"`
	Ip string `json:"ip"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func NewToken(id, ip string) string {
	expiry := time.Now().Add(30 * time.Minute)
	claims := &Claims{Id: id, Ip: ip, StandardClaims: jwt.StandardClaims{ExpiresAt: expiry.Unix()}}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := tkn.SignedString(jwtKey)
	cookie := &http.Cookie{Name: "token", Value: str, Expires: expiry, HttpOnly: false}
	return cookie.String()
}
