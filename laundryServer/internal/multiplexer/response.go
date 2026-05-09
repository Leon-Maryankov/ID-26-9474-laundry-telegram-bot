package multiplexer

import (
	"net/http"
)

// writeJSONResponse writes a JSON response with the specified data and status code.
// It sets the Content-Type header to "application/json".
func writeJSONResponse(writer http.ResponseWriter, data []byte, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(data)
}
