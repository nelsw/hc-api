package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Aggregate struct {
	UserId         string `json:"user_id"`
	SourceIp       string `json:"source_ip"`
	Entry          string `json:"entry"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token,omitempty"`
	TemporaryToken string `json:"temporary_token,omitempty"`
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
	handler       = "hcTokenHandler"
)

// Returns the decoded and unencrypted User Id of the token.
func (e *Aggregate) Id() *string {
	return &e.UserId
}

// This method not only validates entity data but determines the purpose of the entity.
// If our User ID is not empty, we know the purpose is to create a new token.
// Otherwise, we must parse the JWT claims to get the User ID of a JWT token.
func (e *Aggregate) Validate() error {

	// Is the data available valid to create or verify a token?
	if e.SourceIp == "" {
		// Source IP is always required to prevent spoofing.
		return fmt.Errorf("bad token, ip is empty\n")
	} else if e.AccessToken == "" && e.UserId == "" {
		// Only one of these values may be blank.
		return fmt.Errorf("bad token, entry and user id are empty\n")
	}

	// Are we creating a new token?
	if e.UserId != "" {
		return nil // yup
	}

	// Parse entry value into JWT claims and verify complete and valid claim interpretation.
	var claims jwt.StandardClaims
	tokenString := regex.ReplaceAllString(e.AccessToken, `$3`)
	if jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return e.Payload(), nil
	}); err != nil {
		return err // Either the token expired or the signature doesn't match.
	} else if !jwtToken.Valid {
		return fmt.Errorf("bad jwtToken=[%v] claims=[%v]", jwtToken, claims)
	} else if claims.Audience != e.SourceIp {
		return fmt.Errorf("bad ips got=[%s] want=[%s] claims=[%v]", e.SourceIp, claims.Audience, claims)
	} else {
		return nil
	}
}

func (e *Aggregate) Payload() []byte {
	if e.Name() == &temporary {
		return temporaryKey
	} else if e.Name() == &refresh {
		return refreshKey
	} else {
		return accessKey
	}
}

func (*Aggregate) Handler() *string {
	return &handler
}

// Returns the type name of token to be created or verified.
func (e *Aggregate) Name() *string {
	switch e.AccessToken {
	case temporary:
		return &temporary
	case refresh:
		return &refresh
	case access:
		return &access
	default:
		s := regex.ReplaceAllString(e.AccessToken, `$1`)
		return &s
	}
}

func (e *Aggregate) String() string {
	if e.UserId != "" {
		return e.TokenStr()
	} else {
		return e.UserID()
	}
}

func (e *Aggregate) TokenStr() string {
	now := time.Now()
	expires := now
	if e.Name() == &temporary {
		expires.Add(15 * time.Second)
	} else if e.Name() == &refresh {
		expires.Add(30 * time.Minute)
	} else {
		expires.Add(24 * time.Hour)
	}
	jwtToken := &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": signingMethod.Alg(),
			"key": e.Name(),
		},
		Claims: &jwt.StandardClaims{
			Audience:  e.SourceIp,
			ExpiresAt: expires.Unix(),
			Id:        e.UserId,
			IssuedAt:  now.Unix(),
		},
		Method: signingMethod,
	}
	value, _ := jwtToken.SignedString(e.Payload())
	cookie := &http.Cookie{
		Name:     "token",
		Value:    value,
		Expires:  expires,
		HttpOnly: false,
	}
	return cookie.String()
}

func (e *Aggregate) UserID() string {
	// Parse entry value into JWT claims and verify complete and valid claim interpretation.
	var claims jwt.StandardClaims
	tokenString := regex.ReplaceAllString(e.AccessToken, `$3`)
	if jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return e.Payload(), nil
	}); err != nil {
		panic(err)
	} else if !jwtToken.Valid {
		panic(err)
	} else if claims.Audience != e.SourceIp {
		panic(err)
	} else {
		return claims.Id
	}
}
