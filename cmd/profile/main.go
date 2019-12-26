// UserProfile is exactly what it appears to be, a user profile domain model entity.
// It also promotes separation of concerns by decoupling user profile details from the primary user entity. IF
// UserProfile.EmailOld != UserProfile.EmailNew, AND User.Email == UserProfile.EmailOld, THEN we must prompt the user to
// confirm new email address. IF UserProfile.Password1 is not blank AND UserProfile.Password2 is not blank AND valid AND
// UserProfile.Password1 == UserProfile.Password2, then we update the UserPassword entity and return OK.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	. "hc-api/service"
	"os"
)

var table = os.Getenv("USER_PROFILE_TABLE")

type UserProfile struct {
	Id        string   `json:"id"`
	BrandIds  []string `json:"brand_ids"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Phone     string   `json:"phone"`
	// unused
	Password1 string `json:"password_1,omitempty"`
	Password2 string `json:"password_2,omitempty"`
	// deprecated - todo, refactor to qsp.
	Session string `json:"session"`
}

func (up *UserProfile) Unmarshal(s string) error {
	if err := json.Unmarshal([]byte(s), &up); err != nil {
		return err
	} else if up.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		up.Id = id.String()
		return nil
	}
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	body := r.Body
	ip := r.RequestContext.Identity.SourceIP
	session := r.QueryStringParameters["session"]
	fmt.Printf("REQUEST cmd=[%s], ip=[%s], session=[%s], body=[%s]", cmd, ip, session, body)

	switch cmd {

	case "save":
		var p UserProfile
		if err := p.Unmarshal(r.Body); err != nil {
			return BadRequest().Error(err).Build()
		} else if _, err := ValidateSession(p.Session, ip); err != nil {
			return Unauthorized().Error(err).Build()
		} else if err := Put(p, &table); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	case "find":
		var p UserProfile
		id := r.QueryStringParameters["id"]
		if err := FindOne(&table, &id, p); err != nil {
			return NotFound().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	default:
		return BadRequest().Data(r).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
