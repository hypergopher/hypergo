package hyperview

import (
	"encoding/json"
	"net/http"
)

// Envelope represents the structure of an envelope used for encapsulating response data.
type Envelope struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Code    int    `json:"code,omitempty"`
}

// JSONSuccess creates a successful JSON response with the given data and optional headers.
// It uses the Envelope structure to encapsulate the data and set the response status, code, and message.
// It then calls the JSONWithHeaders function to format the JSON response with the specified headers.
// The function returns an error if there is an issue with formatting or writing the response to the writer.
func JSONSuccess(w http.ResponseWriter, data any, headers ...http.Header) error {
	envelope := Envelope{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	}

	return JSONWithHeaders(w, http.StatusOK, envelope, headers...)
}

// JSONSuccessWithStatus creates a JSON response with the specified status code and data.
// It formats the response body as a success envelope and includes optional custom headers.
// It returns an error if writing the response fails.
func JSONSuccessWithStatus(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	envelope := Envelope{
		Status:  "success",
		Code:    status,
		Message: "Success",
		Data:    data,
	}

	return JSONWithHeaders(w, status, envelope, headers...)
}

// JSONFailure builds a JSON response with failure status, message and data.
// It uses the provided http.ResponseWriter to write the JSON response.
// The response code is set by the status parameter.
// The response headers can be passed as optional http.Header arguments.
func JSONFailure(w http.ResponseWriter, data any, message string, status int, headers ...http.Header) error {
	envelope := Envelope{
		Status:  "fail",
		Code:    status,
		Message: message,
		Data:    data,
	}

	return JSONWithHeaders(w, status, envelope, headers...)
}

// JSONError writes an error response in JSON format to the http.ResponseWriter.
// It takes the message and the status code of the error as input parameters.
// Optional headers can be provided to set additional response headers.
//
// Example Usage:
// err := JSONError(w, "Internal Server Error", http.StatusInternalServerError, http.Header{})
//
// Parameters:
// - w: The http.ResponseWriter to write the error response to.
// - message: The error message to be included in the response.
// - status: The HTTP status code of the error.
// - headers: Optional additional response headers.
//
// Returns:
// - error: An error if JSONWithHeaders fails, otherwise nil.
func JSONError(w http.ResponseWriter, message string, status int, headers ...http.Header) error {
	envelope := Envelope{
		Status:  "error",
		Message: message,
		Code:    status,
	}

	return JSONWithHeaders(w, status, envelope, headers...)
}

// JSONRedirect redirects the request to the specified URL and sends a JSON response.
func JSONRedirect(w http.ResponseWriter, r *http.Request, url string, headers ...http.Header) error {
	return JSONWithHeaders(w, http.StatusSeeOther, map[string]string{
		"Redirect": url,
	}, headers...)
}

// JSONWithHeaders serializes the given data to JSON format with specified headers
// and writes it to the provided http.ResponseWriter. It also sets the HTTP status
// code and the Content-Type header to "application/json; charset=UTF-8". If the
// serialization fails, an error is returned. The function accepts optional headers
// that will be applied to the response.
func JSONWithHeaders(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for _, header := range headers {
		for key, value := range header {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = w.Write(js)

	return nil
}
