package entity

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"hc-api/pkg/value"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Token struct {
	value.Authorization
}

// Three types of tokens exist, each with variable expiration times.
var (
	temporary     = "tep"
	refresh       = "ref"
	access        = "acc"
	regex         = regexp.MustCompile(`(.*)(token=)(.*)(;.*)`)
	accessKey     = []byte(os.Getenv("JWT_KEY"))
	refreshKey    = []byte(os.Getenv("REF_JWT_KEY"))
	temporaryKey  = []byte(os.Getenv("TMP_JWT_KEY"))
	signingMethod = jwt.SigningMethodES256
	tokenHandler  = "hcTokenHandler"
)

// This method not only validates entity data but determines the purpose of the entity.
// If our User ID is not empty, we know the purpose is to create a new token.
// Otherwise, we must parse the JWT claims to get the User ID of a JWT token.
func (token *Token) Validate() error {

	// Is the data available valid to create or verify a token?
	if token.SourceIp == "" {
		// Source IP is always required to prevent spoofing.
		return fmt.Errorf("bad token, ip is empty\n")
	} else if token.AccessToken == "" && token.UserId == "" {
		// Only one of these values may be blank.
		return fmt.Errorf("bad token, entry and user id are empty\n")
	}

	// Are we creating a new token?
	if token.UserId != "" {
		return nil // yup
	}

	// Parse entry value into JWT claims and verify complete and valid claim interpretation.
	var claims jwt.StandardClaims
	tokenString := regex.ReplaceAllString(token.AccessToken, `$3`)
	if jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(tkn *jwt.Token) (i interface{}, err error) {
		return token.Key(), nil
	}); err != nil {
		return err // Either the token expired or the signature doesn't match.
	} else if !jwtToken.Valid {
		return fmt.Errorf("bad jwtToken=[%v] claims=[%v]", jwtToken, claims)
	} else if claims.Audience != token.SourceIp {
		return fmt.Errorf("bad ips got=[%s] want=[%s] claims=[%v]", token.SourceIp, claims.Audience, claims)
	} else {
		return nil
	}
}

func (token *Token) Key() []byte {
	if token.Type() == &temporary {
		return temporaryKey
	} else if token.Type() == &refresh {
		return refreshKey
	} else {
		return accessKey
	}
}

func (*Token) Function() *string {
	return &tokenHandler
}

func (token *Token) Payload() []byte {
	if token.UserId != "" {
		return []byte(token.String())
	} else {
		return []byte(token.Id())
	}
}

// Returns the type name of token to be created or verified.
func (token *Token) Type() *string {
	switch token.AccessToken {
	case temporary:
		return &temporary
	case refresh:
		return &refresh
	case access:
		return &access
	default:
		s := regex.ReplaceAllString(token.AccessToken, `$1`)
		return &s
	}
}

func (token *Token) String() string {
	now := time.Now()
	expires := now
	if token.Type() == &temporary {
		expires.Add(15 * time.Second)
	} else if token.Type() == &refresh {
		expires.Add(30 * time.Minute)
	} else {
		expires.Add(24 * time.Hour)
	}
	jwtToken := &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": signingMethod.Alg(),
			"key": token.Type(),
		},
		Claims: &jwt.StandardClaims{
			Audience:  token.SourceIp,
			ExpiresAt: expires.Unix(),
			Id:        token.UserId,
			IssuedAt:  now.Unix(),
		},
		Method: signingMethod,
	}
	value, _ := jwtToken.SignedString(token.Key())
	cookie := &http.Cookie{
		Name:     "token",
		Value:    value,
		Expires:  expires,
		HttpOnly: false,
	}
	return cookie.String()
}

func (token *Token) Id() string {
	// Parse entry value into JWT claims and verify complete and valid claim interpretation.
	var claims jwt.StandardClaims
	str := regex.ReplaceAllString(token.AccessToken, `$3`)
	if jwtToken, err := jwt.ParseWithClaims(str, &claims, func(tkn *jwt.Token) (i interface{}, err error) {
		return token.Key(), nil
	}); err != nil {
		panic(err)
	} else if !jwtToken.Valid {
		panic(err)
	} else if claims.Audience != token.SourceIp {
		panic(err)
	} else {
		return claims.Id
	}
}
