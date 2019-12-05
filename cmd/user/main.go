package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	response "github.com/nelsw/hc-util/aws"
	"golang.org/x/crypto/bcrypt"
	"hc-api/service"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var userTable = os.Getenv("USER_TABLE")
var userPasswordTable = os.Getenv("USER_PASSWORD_TABLE")
var userEmailTable = os.Getenv("USER_EMAIL_TABLE")

// No email regex is perfect, but this one is close.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Primary user object for the domain, visible to client and server. With the exception of the Email field, which
// represents its plain text value, each property is a reference to a unique ID, or collection of unique ID's.
type User struct {
	Id         string   `json:"id"`
	PasswordId string   `json:"password_id"`
	ProfileId  string   `json:"profile_id"`
	AddressIds []string `json:"address_ids"`
	ProductIds []string `json:"product_ids"`
	OrderIds   []string `json:"order_ids"`
	SaleIds    []string `json:"sale_ids"`
	Session    string   `json:"session"`
}

// Used for login and registration use cases.
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Data structure for securely associating user entities with their email address.
type UserEmail struct {
	Id     string `json:"id"` // email address
	UserId string `json:"user_id"`
}

// Data structure for persisting and retrieving a users (encrypted) password. The user entity maintains a 1-1 FK
// relationship with the UserPassword entity for referential integrity, where user.PasswordId == userPassword.Id.
type UserPassword struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

// Used to update an existing user item in an Amazon DynamoDB table.
// SET - modify or add item attributes
// REMOVE - delete attributes from an item
// ADD - update numbers and sets
// DELETE - remove elements from a set
type UserUpdate struct {
	Val        []string `json:"val"`
	Expression string   `json:"expression"`
	Session    string   `json:"session"`
}

// Validates the UserCredentials entity by confirming that both the email and password values are valid.
// Allows email addresses with third party domains and any extension.
func (uc *UserCredentials) Validate() error {
	if emailRegex.MatchString(uc.Email) == false {
		return fmt.Errorf("bad email [%s]", uc.Email)
	} else if err := IsPasswordValid(uc.Password); err != nil {
		return err
	} else {
		return nil
	}
}

func (uc *UserCredentials) Unmarshal(s string) error {
	if err := json.Unmarshal([]byte(s), &uc); err != nil {
		return err
	} else if err := uc.Validate(); err != nil {
		return err
	} else {
		uc.Email = strings.ToLower(uc.Email)
		return nil
	}
}

// Validates the UserPassword entity by confirming that both the password and id values are valid.
func (up *UserPassword) Validate() error {
	return IsPasswordValid(up.Password)
}

// The following is an adaptation of https://stackoverflow.com/a/25840157
func IsPasswordValid(s string) error {
	var number, upper, special bool
	length := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
			length++
		case unicode.IsUpper(c):
			upper = true
			length++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			length++
		case unicode.IsLetter(c) || c == ' ':
			length++
		default:
			// do not increment length for unrecognized characters
		}
	}
	if length < 8 || length > 24 {
		return fmt.Errorf("bad password, must contain 8-24 characters")
	} else if number == false {
		return fmt.Errorf("bad password, must contain at least 1 number")
	} else if upper == false {
		return fmt.Errorf("bad password, must contain at least 1 uppercase letter")
	} else if special == false {
		return fmt.Errorf("bad password, must contain at least 1 special character")
	} else {
		return nil
	}
}

// Finds user by id (PK).
func findUserById(s *string) (user *User, err error) {
	if result, err := service.Get(&userTable, s); err == nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	}
	return user, err
}

// Returns the user password entity by providing the user password primary key.
func findUserPasswordById(s *string) (up *UserPassword, err error) {
	if result, err := service.Get(&userPasswordTable, s); err == nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &up)
	}
	return up, err
}

// Returns the user email entity by providing the user email primary key.
func findUserEmailById(s *string) (ue *UserEmail, err error) {
	if result, err := service.Get(&userEmailTable, s); err == nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &ue)
	}
	return ue, err
}

// Updates the specified attributes of a user entity.
func updateUser(k, e *string, v *[]string) error {
	return service.Update(&dynamodb.UpdateItemInput{
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        &userTable,
		Key:              map[string]*dynamodb.AttributeValue{"id": {S: k}},
		UpdateExpression: e,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				SS: aws.StringSlice(*v),
			},
		},
	})
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	fmt.Printf("REQUEST [%s]: [%v]", cmd, r)

	switch cmd {

	case "register":
		// todo - create User, UserPassword, and UserProfile entities ... also verify email address.
		return response.New().Code(http.StatusNotImplemented).Build()

	case "login":
		var uc UserCredentials
		if err := uc.Unmarshal(r.Body); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if ue, err := findUserEmailById(&uc.Email); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if u, err := findUserById(&ue.UserId); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if up, err := findUserPasswordById(&u.PasswordId); err != nil {
			return response.New().Code(http.StatusNotFound).Text(err.Error()).Build()
		} else if err := bcrypt.CompareHashAndPassword([]byte(up.Password), []byte(uc.Password)); err != nil {
			return response.New().Code(http.StatusUnauthorized).Build()
		} else if cookie, err := service.NewCookie(u.Id, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusInternalServerError).Build()
		} else {
			u.Session = cookie
			u.OrderIds = nil
			u.SaleIds = nil
			u.Id = ""
			u.PasswordId = ""
			return response.New().Code(http.StatusOK).Data(&u).Build()
		}

	case "update":
		var uu UserUpdate
		if err := json.Unmarshal([]byte(r.Body), &uu); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if id, err := service.ValidateSession(uu.Session, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := updateUser(&id, &uu.Expression, &uu.Val); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Text(string([]byte(`{"success":""}`))).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
