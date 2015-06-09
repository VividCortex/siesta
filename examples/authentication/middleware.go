package main

import (
	"github.com/VividCortex/siesta"

	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type apiResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func requestIdentifier(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	requestID := fmt.Sprintf("%x", rand.Int())
	c.Set("request-id", requestID)
	log.Printf("[Req %s] %s %s", requestID, r.Method, r.URL)
}

func authenticator(c siesta.Context, w http.ResponseWriter, r *http.Request,
	quit func()) {
	requestID := c.Get("request-id").(string)
	db := c.Get("db").(*state)

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
		c.Set("user", user)
	} else {
		log.Printf("[Req %s] Did not provide a token", requestID)

		c.Set("error", "token required")
		c.Set("status-code", http.StatusUnauthorized)

		quit()
		return
	}
}

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

func jsonResponseWriter(c siesta.Context, w http.ResponseWriter, r *http.Request,
	quit func()) {
	// Set the request ID header.
	if requestID := c.Get("request-id"); requestID != nil {
		w.Header().Set("X-Request-ID", requestID.(string))
	}

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
		w.Header().Set("Content-Type", "application/json")
		// We'll encode it as JSON without knowing
		// what it exactly is.
		enc.Encode(response)
	}

	// We're at the end of the middleware chain, so quit.
	quit()
}
