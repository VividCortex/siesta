package main

import (
	"github.com/VividCortex/siesta"

	"log"
	"net/http"
)

func getResource(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	requestID := c.Get("request-id").(string)
	db := c.Get("db").(*state)

	user := c.Get("user").(string)
	var params siesta.Params
	resourceID := params.Int("resourceID", -1, "Resource identifier")
	err := params.Parse(r.Form)

	if err != nil {
		log.Printf("[Req %s] %v", requestID, err)
		c.Set("error", err.Error())
		c.Set("status-code", http.StatusBadRequest)
		return
	}

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
