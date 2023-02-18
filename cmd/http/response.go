package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Write HTTP Error out to http.ResponseWriter
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	var errorCode string
	switch statusCode {
	case 400:
		errorCode = "InvalidParameter."
	case 404:
		errorCode = "NotFound or caller don't have access."
	case 401:
		errorCode = "Unauthorized"
	case 403:
		errorCode = "Unauthorized"
	case 409:
		errorCode = "Conflict with resource"
	case 500:
		errorCode = "InternalServerError"
	}
	log.Warnf("Received error from call %s", err)
	w.Write(SerializeError(statusCode, errorCode))

}
