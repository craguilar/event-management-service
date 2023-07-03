package main

import (
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

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func (h *LambaHandler) Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf(
		"Received request %s %s %s",
		request.HTTPMethod,
		request.Path,
		request.Body,
	)
	response, err := h.adapter.Proxy(*core.NewSwitchableAPIGatewayRequestV1(&request))
	if err != nil {
		log.Error(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "InternalServerError",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
	return *response.Version1(), nil
}
