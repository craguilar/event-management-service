package http

import (
	"encoding/json"
)

// Internal error object model.
type ModelError struct {

	// A short error code that defines the error, meant for programmatic parsing.
	Code int `json:"code"`

	// A human-readable error string.
	Message string `json:"message"`
}

func SerializeData(data interface{}) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return result
}

func SerializeError(httpCode int, message string) []byte {
	result, error := json.Marshal(ModelError{Code: httpCode, Message: message})
	if error != nil {
		return []byte{}
	}
	return result
}
