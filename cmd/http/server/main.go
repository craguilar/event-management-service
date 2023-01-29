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

	// Create services and provide it to handler
	event := mock.NewEventService()
	guest := mock.NewGuestService(event)
	handler := appHttp.NewServiceHandler(event, guest)
	// Router config
	router := appHttp.NewRouter(handler)

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cmd.GetConfig("PORT")), router))
}
