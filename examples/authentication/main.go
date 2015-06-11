package main

import (
	"github.com/VividCortex/siesta"

	"log"
	"net/http"
)

func main() {
	// Create a new service rooted at /.
	service := siesta.NewService("/")

	// requestIdentifier assigns an ID to every request
	// and adds it to the context for that request.
	// This is useful for logging.
	service.AddPre(requestIdentifier)

	// Add access to the state via the context in every handler.
	service.AddPre(func(c siesta.Context, w http.ResponseWriter, r *http.Request) {
		c.Set("db", state)
	})

	// We'll add the authenticator middleware to the "pre" chain.
	// It will ensure that every request has a valid token.
	service.AddPre(authenticator)

	// Response generation
	service.AddPost(responseGenerator)
	service.AddPost(responseWriter)

	// Custom 404 handler
	service.SetNotFound(func(c siesta.Context, w http.ResponseWriter, r *http.Request) {
		c.Set("status-code", http.StatusNotFound)
		c.Set("error", "not found")
	})

	// Routes
	service.Route("GET", "/resources/:resourceID", "Retrieves a resource",
		getResource)

	log.Println("Listening on :8080")
	panic(http.ListenAndServe(":8080", service))
}
