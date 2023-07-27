package app

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	CharSet = "UTF-8"
)

type EmailConfig struct {
	SesService *ses.SES
}

type EmailNotificationService struct {
	config *EmailConfig
}

func NewEmailNotificationService(config *EmailConfig) *EmailNotificationService {
	return &EmailNotificationService{
		config: config,
	}
}

// From https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/ses-example-send-email.html
func (e *EmailNotificationService) SendEmailNotification(recipient, subject string, body string) error {

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		// TODO: Remove hardcoding.
		Source: aws.String("app@notifications.cmymesh.com"),
	}

	// Attempt to send the email.
	result, err := e.config.SesService.SendEmail(input)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("Email Sent to address %s %v", recipient, result)
	return nil
}
