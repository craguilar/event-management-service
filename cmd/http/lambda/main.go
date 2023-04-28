package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	appHttp "github.com/craguilar/event-management-service/cmd/http"
	"github.com/craguilar/event-management-service/internal/app/dynamo"
)

var db *dynamo.DBConfig

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
	handler := appHttp.NewServiceHandler(event, guest, task, expense)
	// Router and Lambda Handler
	router := appHttp.NewRouter(handler)
	lambdHandler := NewLambaHandler(router)
	// Start lambda
	lambda.Start(lambdHandler.Handler)
}
