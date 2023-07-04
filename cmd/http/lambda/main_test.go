package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	appHttp "github.com/craguilar/event-management-service/cmd/http"
	"github.com/craguilar/event-management-service/internal/app/mock"
	"github.com/stretchr/testify/assert"
)

func TesGettHandler(t *testing.T) {

	request := events.APIGatewayProxyRequest{}
	request.Path = "/"
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "404 page not found",
	}

	handler := createMockHandler()
	response, err := handler.HandleHttp(request)

	// assert.Equal(t, response.Headers, expectedResponse.Headers)
	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)

}

func TesOptionstHandler(t *testing.T) {

	request := events.APIGatewayProxyRequest{}
	request.Path = "/20230125"
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "",
	}

	handler := createMockHandler()
	response, err := handler.HandleHttp(request)

	assert.Contains(t, expectedResponse.StatusCode, response.StatusCode)
	assert.Equal(t, err, nil)

}

func createMockHandler() *LambaHandler {
	// Create service and provide it to handler
	event := mock.NewEventService()
	guest := mock.NewGuestService(event)
	task := &mock.TaskService{}
	expense := &mock.ExpenseService{}
	handler := appHttp.NewServiceHandler(event, guest, task, expense)
	router := appHttp.NewRouter(handler)
	return NewLambaHandler(router)
}
