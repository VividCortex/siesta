package main

import (
	"github.com/VividCortex/siesta"

	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

// apiResponse defines the structure of the responses.
type apiResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// requestIdentifier generates a request ID and sets the "request-id"
// key in the context. It also logs the request ID and the requested URL.
func requestIdentifier(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	requestID := fmt.Sprintf("%x", rand.Int())
	c.Set("request-id", requestID)
	log.Printf("[Req %s] %s %s", requestID, r.Method, r.URL)
}

// authenticator reads the username from the HTTP basic authentication header
// and validates the token. It sets the "user" key in the context to the
// user associated with the token.
func authenticator(c siesta.Context, w http.ResponseWriter, r *http.Request,
	quit func()) {
	// Context variables
	requestID := c.Get("request-id").(string)
	db := c.Get("db").(*DB)

	// Check for a token in the HTTP basic authentication username field.
	token, _, ok := r.BasicAuth()
	if ok {
		user, err := db.validateToken(token)
		if err != nil {
			log.Printf("[Req %s] Did not provide a valid token", requestID)
			c.Set("status-code", http.StatusUnauthorized)
			c.Set("error", "invalid token")
			quit()
			return
		}

		log.Printf("[Req %s] Provided a token for: %s", requestID, user)

		// Add the user to the context.
		c.Set("user", user)
	} else {
		log.Printf("[Req %s] Did not provide a token", requestID)

		c.Set("error", "token required")
		c.Set("status-code", http.StatusUnauthorized)

		// Exit the chain here.
		quit()
		return
	}
}

// responseGenerator converts response and/or error data passed through the
// context into a structured response.
func responseGenerator(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	response := apiResponse{}

	if data := c.Get("data"); data != nil {
		response.Data = data
	}

	if err := c.Get("error"); err != nil {
		response.Error = err.(string)
	}

	c.Set("response", response)
}

// responseWriter sets the proper headers and status code, and
// writes a JSON-encoded response to the client.
func responseWriter(c siesta.Context, w http.ResponseWriter, r *http.Request,
	quit func()) {
	// Set the request ID header.
	if requestID := c.Get("request-id"); requestID != nil {
		w.Header().Set("X-Request-ID", requestID.(string))
	}

	// Set the content type.
	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)

	// If we have a status code set in the context,
	// send that in the header.
	//
	// Go defaults to 200 OK.
	statusCode := c.Get("status-code")
	if statusCode != nil {
		statusCodeInt := statusCode.(int)
		w.WriteHeader(statusCodeInt)
	}

	// Check to see if we have some sort of response.
	response := c.Get("response")
	if response != nil {
		// We'll encode it as JSON without knowing
		// what it exactly is.
		enc.Encode(response)
	}

	// We're at the end of the middleware chain, so quit.
	quit()
}
