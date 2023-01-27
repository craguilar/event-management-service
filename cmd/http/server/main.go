package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/craguilar/event-management-service/cmd"
	appHttp "github.com/craguilar/event-management-service/cmd/http"
	"github.com/craguilar/event-management-service/internal/app/mock"
)

func main() {
	log.Printf("Server started on port %s", cmd.GetConfig("PORT"))

	// Create car service and provide it to handler
	eventService := mock.NewEventService()
	handler := appHttp.NewServiceHandler(eventService)
	router := appHttp.NewRouter(handler)

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cmd.GetConfig("PORT")), router))
}
