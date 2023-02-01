package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Write HTTP Error out to http.ResponseWriter
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	var errorCode string
	if statusCode == 400 {
		errorCode = "InvalidParameter."
	} else if statusCode == 404 {
		errorCode = "NotFound or caller don't have access."
	} else if statusCode == 401 || statusCode == 403 {
		errorCode = "Unauthorized"
	} else if statusCode == 409 {
		errorCode = "Conflict with resource"
	} else if statusCode == 500 {
		errorCode = "InternalServerError"
	}
	log.Warnf("Received error from call %s", err)
	w.Write(SerializeError(statusCode, errorCode))

}
