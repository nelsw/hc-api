package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Entity struct {
	Value
	Error
}

type Error struct {
	Msg string `json:"errorMessage"`
	Typ string `json:"errorType"`
}

type Value struct {
	SourceId string   `json:"source_id"`
	SourceIp string   `json:"source_ip"`
	JwtSlice []string `json:"jwt_slice"` // *-token=*.*.*; Expires=Wed, 29 Jan 2020 04:07:30 GMT
	Subject  string   `json:"subject"`
	Msg      string   `json:"errorMessage"`
	Typ      string   `json:"errorType"`
}

type Jwt struct {
	Name     string
	Duration time.Duration
	JwtKey   []byte
}

var (
	signingMethod      = jwt.SigningMethodHS256
	functionName       = "hcTokenHandler"
	regex              = regexp.MustCompile(`(.*)(token=)(.*)(;.*)`)
	access             Jwt
	refresh            Jwt
	critical           Jwt
	InvalidToken       = fmt.Errorf("bad token\n")
	ErrBadIpData       = fmt.Errorf("bad token, source ip is empty\n")
	ErrBadJwtToken     = fmt.Errorf("bad token, invalid segments or expired\n")
	ErrBadCookieData   = fmt.Errorf("bad token, cookies and user identity are empty\n")
	ErrBadIpComparison = fmt.Errorf("bad token, unable to match source and target ip validation\n")
)

func init() {
	access = Jwt{"acc", time.Hour * 24, []byte(os.Getenv("ACC_JWT_KEY"))}
	refresh = Jwt{"ref", time.Minute * 30, []byte(os.Getenv("REF_JWT_KEY"))}
	critical = Jwt{"crt", time.Second * 15, []byte(os.Getenv("CRT_JWT_KEY"))}
}

func (v *Jwt) KeyFunc() jwt.Keyfunc {
	return func(tkn *jwt.Token) (i interface{}, err error) {
		return v.JwtKey, nil
	}
}

func (e *Entity) Function() string {
	return functionName
}

// Validations common for creating and verifying a token.
func (e *Entity) Validate() error {
	if e.SourceIp == "" {
		return ErrBadIpData // Source IP is always required to prevent spoofing.
	} else if (e.JwtSlice == nil || len(e.JwtSlice) == 0) && e.SourceId == "" {
		return ErrBadCookieData // Only one of these values may be blank.
	} else {
		return nil
	}
}

// Validate each existing token query.
func (e *Entity) Authenticate() error {

	var str string
	var keyFunc jwt.Keyfunc
	var claims jwt.StandardClaims

	for _, s := range e.JwtSlice {

		switch regex.ReplaceAllString(s, `$1`) {
		case critical.Name:
			keyFunc = critical.KeyFunc()
		case refresh.Name:
			keyFunc = refresh.KeyFunc()
		default:
			keyFunc = access.KeyFunc()
		}

		str = regex.ReplaceAllString(s, `$3`)

		if jwtToken, err := jwt.ParseWithClaims(str, &claims, keyFunc); err != nil {
			return err // Either the token expired or signatures do not match.
		} else if !jwtToken.Valid {
			return ErrBadJwtToken
		} else if claims.Audience != e.SourceIp {
			return ErrBadIpComparison
		}
	}
	e.SourceId = claims.Id
	return nil
}

func (e *Entity) cookie() Jwt {
	if len(e.JwtSlice) > 1 {
		return critical
	} else if len(e.JwtSlice) < 1 {
		return access
	} else {
		return refresh
	}
}

func (e *Entity) Authorize() error {

	var tknType = e.cookie()

	now := time.Now()
	expires := now.Add(tknType.Duration)
	jwtToken := &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": signingMethod.Alg(),
		},
		Claims: jwt.StandardClaims{
			Audience:  e.SourceIp,
			ExpiresAt: expires.Unix(),
			Id:        e.SourceId,
			IssuedAt:  now.Unix(),
			Subject:   e.Subject,
		},
		Method: signingMethod,
	}

	val, err := jwtToken.SignedString(tknType.JwtKey)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:    tknType.Name + "-token",
		Value:   val,
		Expires: expires,
	}

	e.JwtSlice = append(e.JwtSlice, cookie.String())

	return nil
}
