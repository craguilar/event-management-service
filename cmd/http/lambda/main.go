package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ses"
	appHttp "github.com/craguilar/event-management-service/cmd/http"
	"github.com/craguilar/event-management-service/internal/app"
	"github.com/craguilar/event-management-service/internal/app/dynamo"
)

var db *dynamo.DBConfig
var emailConfig *app.EmailConfig

func init() {
	log.Println("Initializing DynamoDB lambda")
	if db == nil {

		region := os.Getenv("AWS_REGION")
		awsSession, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			log.Fatalf("Error found %s", err)
			return
		}
		db = dynamo.InitDb(dynamodb.New(awsSession), "events")
	}
	if emailConfig == nil {
		region := os.Getenv("AWS_REGION")
		sesSession, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			log.Fatalf("Error found %s", err)
			return
		}
		emailConfig = &app.EmailConfig{
			SesService: ses.New(sesSession),
		}
	}
	log.Println("Initialized lambda ", db.DbService.Endpoint)
}

func main() {
	log.Printf("Lambda started")
	// Authorization
	authorize := dynamo.NewAuthorizationService(db)
	// Create services and provide it to handler
	event := dynamo.NewEventService(db, authorize)
	guest := dynamo.NewGuestService(db, authorize)
	task := dynamo.NewTaskService(db)
	expense := dynamo.NewExpenseService(db)
	notification := app.NewEmailNotificationService(emailConfig)
	actions := dynamo.NewEventActionsService(db, event, task, notification)
	handler := appHttp.NewServiceHandler(event, actions, guest, task, expense)
	// Router and Lambda Handler
	router := appHttp.NewRouter(handler)
	lambdHandler := NewLambaHandler(router)
	// Start lambda
	lambda.Start(lambdHandler.Handler)
}
