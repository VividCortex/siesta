package main

import (
	"github.com/VividCortex/siesta"

	"log"
	"net/http"
)

// getResource is the function that handles the GET /resources/:resourceID route.
func getResource(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	// Context variables
	requestID := c.Get("request-id").(string)
	db := c.Get("db").(*DB)
	user := c.Get("user").(string)

	// Check parameters
	var params siesta.Params
	resourceID := params.Int("resourceID", -1, "Resource identifier")
	err := params.Parse(r.Form)
	if err != nil {
		log.Printf("[Req %s] %v", requestID, err)
		c.Set("error", err.Error())
		c.Set("status-code", http.StatusBadRequest)
		return
	}

	// Make sure we have a valid resource ID.
	if *resourceID == -1 {
		c.Set("error", "invalid or missing resource ID")
		c.Set("status-code", http.StatusBadRequest)
		return
	}

	resource, err := db.resource(user, *resourceID)
	if err != nil {
		c.Set("status-code", http.StatusNotFound)
		c.Set("error", "not found")
		return
	}

	c.Set("data", resource)
}
