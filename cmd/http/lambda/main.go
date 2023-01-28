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
	log.Println("Initializing lambda")
	if db == nil {

		region := os.Getenv("AWS_REGION")
		awsSession, err := session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		if err != nil {
			log.Fatalf("Error found %s", err)
			return
		}
		db = dynamo.InitDb(dynamodb.New(awsSession), "cars")
	}
	log.Println("Initialized lambda ", db.DbService.Endpoint)
}

func main() {
	log.Printf("Lambda started")
	// Create car service and provide it to handler
	service := dynamo.NewEventService(db)
	serviceHandler := appHttp.NewServiceHandler(service)
	router := appHttp.NewRouter(serviceHandler)
	handler := NewLambaHandler(router)
	// Start lambda
	lambda.Start(handler.Handler)
}
