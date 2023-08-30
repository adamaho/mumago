package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorCode string

const (
	DatabaseError        ErrorCode = "DATABASE_ERROR"
	PatchError           ErrorCode = "PATCH_ERROR"
	MarshalError         ErrorCode = "MARSHAL_ERROR"
	InvalidRequestBody   ErrorCode = "INVALID_REQUEST_BODY"
	InvalidRequestParams ErrorCode = "INVALID_REQUEST_PARAMS"
)

// Converts an ErrorCode into a status code
func StatusCode(e ErrorCode) int {
	switch e {
	case DatabaseError:
		return http.StatusInternalServerError
	case MarshalError:
		return http.StatusInternalServerError
	case InvalidRequestBody:
		return http.StatusBadRequest
	case InvalidRequestParams:
		return http.StatusBadRequest
	case PatchError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Converts an ErrorCode into a human readable string
func Message(e ErrorCode) string {
	switch e {
	case DatabaseError:
		return "An error occured while trying to execute a database operation"
	case MarshalError:
		return "An error occured while trying marshal json"
	case InvalidRequestBody:
		return "Invalid request body was provided"
	case InvalidRequestParams:
		return "Invalid request params were provided"
	case PatchError:
		return "Failed to generate patch"
	default:
		return "An internal error occured"
	}
}

// Representation of the response body of an error
type httpErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Info    string    `json:"info"`
}

// Returns an http error response
func HttpError(w http.ResponseWriter, c ErrorCode, info string) {
	response, err := json.Marshal(httpErrorResponse{Code: c, Message: Message(c), Info: info})

	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Oops. If you are seeing this, something went really wrong."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
