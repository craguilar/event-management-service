package main

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
)

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
	// Try parsing into a APIGatewayProxyRequest
	http := events.APIGatewayProxyRequest{}
	if err := json.Unmarshal(request, &http); err == nil {

		return h.HandleHttp(http)
	}
	// Non handled code
	return getErrorResponse(errors.New("type not enabled"))
}

// Handle HTTP , returns an Amazon API Gateway response object to AWS Lambda
func (h *LambaHandler) HandleHttp(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf(
		"Received request %s %s",
		request.HTTPMethod,
		request.Path,
	)
	response, err := h.adapter.Proxy(*core.NewSwitchableAPIGatewayRequestV1(&request))
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
