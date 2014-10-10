// utils.go contains common code that is not necessary to use with siesta,
// but is often helpful.

package util

import (
	"github.com/VividCortex/siesta"

	"encoding/json"
	"log"
	"net/http"
)

// JsonResponseWriter is a common "post" middleware that writes the
// response as JSON and places the status code in the header. It looks
// for the status code in context[statusCodeStr] (must be castable to int).
// And it looks for response in context[response].
func JsonResponseWriter(statusCodeStr, responseStr string) func(c siesta.Context, w http.ResponseWriter, r *http.Request, quit func()) {
	return func(c siesta.Context, w http.ResponseWriter, r *http.Request, quit func()) {
		enc := json.NewEncoder(w)

		// If we have a status code set in the context,
		// send that in the header.
		//
		// Go defaults to 200 OK.
		statusCode := c.Get(statusCodeStr)
		if statusCode != nil {
			statusCodeInt := statusCode.(int)
			w.WriteHeader(statusCodeInt)
		}

		// Check to see if we have some sort of response.
		response := c.Get(responseStr)
		if response != nil {
			// We'll encode it as JSON without knowing
			// what it exactly is.
			err := enc.Encode(response)
			if err != nil {
				log.Println("couldn't encode response:", err)
			}
		}
	}
}
