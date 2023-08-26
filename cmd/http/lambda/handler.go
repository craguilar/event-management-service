package main

import (
	"encoding/json"
	"errors"

	cmdHttp "github.com/craguilar/event-management-service/cmd/http"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/gorilla/mux"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
)

type ScheduledRequest struct {
	Type string `json:"type"`
}

type LambaHandler struct {
	adapter *gorillamux.GorillaMuxAdapter
}

func NewLambaHandler(r *mux.Router) *LambaHandler {
	gorillaAdapter := gorillamux.New(r)
	return &LambaHandler{
		adapter: gorillaAdapter,
	}
}

// Handler is executed by AWS Lambda in the main function. We are making this handler
// generic so it can receive muliple type of events but always send back an APIGatewayProxyResponse
func (h *LambaHandler) Handler(raw map[string]interface{}) (events.APIGatewayProxyResponse, error) {

	request, err := json.Marshal(raw)
	// Try to marshall the map into a json encoding represented as a byte[]
	if err != nil {
		// do error check
		log.Error(err)
		return getErrorResponse(err)
	}
	http := events.APIGatewayProxyRequest{}
	isScheduled := slices.Contains[string](cmdHttp.TASKS_PATH, http.Path)
	// Try parsing into a APIGatewayProxyRequest
	if err := json.Unmarshal(request, &http); err == nil && !isScheduled && http.HTTPMethod != "" {
		return h.HandleHttp(http)
	}
	// Filter only allowed tasks
	scheduled := ScheduledRequest{}
	if err := json.Unmarshal(request, &scheduled); err == nil && isScheduled && scheduled.Type != "" {
		return h.InterceptScheduled(scheduled)
	}
	// Non handled code
	return getErrorResponse(errors.New("type not enabled"))
}

func (h *LambaHandler) InterceptScheduled(scheduled ScheduledRequest) (events.APIGatewayProxyResponse, error) {
	if scheduled.Type != "PENDING_TASKS" {
		return getErrorResponse(errors.New("invalid scheduled type"))
	}
	// TODO: Could we do this in a better way ?
	request := events.APIGatewayProxyRequest{
		Path:       cmdHttp.TASK_NOTIFICATION,
		HTTPMethod: "POST",
		Headers:    map[string]string{"Authorization": "Bearer dummy"},
	}
	return h.HandleHttp(request)
}

// Handle HTTP , returns an Amazon API Gateway response object to AWS Lambda
func (h *LambaHandler) HandleHttp(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf(
		"HTTP received request %s %s",
		request.HTTPMethod,
		request.Path,
	)
	response, err := h.adapter.Proxy(*core.NewSwitchableAPIGatewayRequestV1(&request))
	log.Printf("Response %v %v ", response.Version1().StatusCode, err)
	if err != nil {
		return getErrorResponse(err)
	}
	return *response.Version1(), nil
}

// Take and convert it into an APIGatewayProxyResponse with status 500
func getErrorResponse(err error) (events.APIGatewayProxyResponse, error) {
	log.Error(err)
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       "InternalServerError :" + err.Error(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, err
}
