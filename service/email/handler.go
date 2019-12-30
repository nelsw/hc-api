package email

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
	"io/ioutil"
	"log"
	"os"
)

type Email struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Code     string `json:"code"`
	Template string `json:"template"`
}

const bucket = "hempconduit"
const key = "template/password-reset.html"

var mgr *s3manager.Downloader
var svc *ses.SES

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}); err != nil {
		log.Printf("Failed to connect to AWS: %s", err.Error())
	} else {
		svc = ses.New(sess)
		mgr = s3manager.NewDownloader(sess)
	}
}

func SendEmail(s string) error {

	f, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	// Write the contents of S3 Object to the file
	if _, err = mgr.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}

	var e Email
	if err := json.Unmarshal([]byte(s), &e); err != nil {
		return err
	} else if b, err := ioutil.ReadFile(f.Name()); err != nil {
		return err
	} else if _, err := svc.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: []byte("Reply-To: noreply@hempconduit.com" +
				"\r\n" + "From: noreply@hempconduit.com" +
				"\r\n" + "To: " + e.To +
				"\r\n" + "Cc: " +
				"\r\n" + "Bcc: " +
				"\r\n" + "Subject: " + e.Subject +
				"\r\n" + "MIME-Version: 1.0" +
				"\r\n" + "Content-Type: text/html; charset=\"utf-8\"\r\n" +
				"\r\n" + string(b) + "\r\n"),
		},
	}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
			return aerr
		} else {
			return err
		}
	} else if err := f.Close(); err != nil {
		return err
	} else {
		return nil
	}
}
