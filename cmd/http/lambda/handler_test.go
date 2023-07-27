package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/craguilar/event-management-service/cmd/http"
	appHttp "github.com/craguilar/event-management-service/cmd/http"
	"github.com/craguilar/event-management-service/internal/app/dynamo"
	"github.com/craguilar/event-management-service/internal/app/mock"
)

func TestPathRouterForScheduledActions(t *testing.T) {

	// Mock services, Router and Lambda Handler
	event := mock.NewEventService()
	guest := mock.NewGuestService(event)
	task := &mock.TaskService{}
	expense := dynamo.NewExpenseService(db)
	action := mock.NewEventActionsService(event, task)
	handler := appHttp.NewServiceHandler(event, action, guest, task, expense)
	router := appHttp.NewRouter(handler)
	lambdHandler := NewLambaHandler(router)
	// Prepare
	request := events.APIGatewayProxyRequest{
		Path:       http.BASE_PATH + "/events/actions/notifyPendingTasks",
		HTTPMethod: "POST",
	}
	response, err := lambdHandler.HandleHttp(request)
	if err != nil {
		t.Fail()
	}
	if response.StatusCode != 200 {
		t.Fail()
	}
}
