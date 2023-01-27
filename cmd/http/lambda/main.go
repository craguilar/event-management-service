package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

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
	/*
		carService := dynamo.NewCarService(db)
		carHandler := appHttp.NewCarServiceHandler(carService)
		router := appHttp.NewRouter(carHandler)
		handler := NewLambaHandler(router)
		// Start lambda
		lambda.Start(handler.Handler)
	*/
}
