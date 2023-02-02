package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/craguilar/event-management-service/cmd"
	appHttp "github.com/craguilar/event-management-service/cmd/http"
	"github.com/craguilar/event-management-service/internal/app/dynamo"
)

// TODO: To refactor this db assignment
var db *dynamo.DBConfig

func init() {
	log.Println("Initializing lambda")
	if db == nil {
		db = dynamo.InitLocalDb("http://localhost:8000", "events")
	}
	log.Println("Initialized lambda ", db.DbService.Endpoint)
}

func main() {
	log.Printf("Server started on port %s", cmd.GetConfig("PORT"))

	// Create services and provide it to handler
	event := dynamo.NewEventService(db)
	guest := dynamo.NewGuestService(db)
	task := dynamo.NewTaskService(db)
	expense := dynamo.NewExpenseService(db)
	handler := appHttp.NewServiceHandler(event, guest, task, expense)
	// Router config
	router := appHttp.NewRouter(handler)

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cmd.GetConfig("PORT")), router))
}
