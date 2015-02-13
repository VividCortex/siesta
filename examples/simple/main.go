package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VividCortex/siesta"
)

func main() {
	// Create a new Service rooted at "/"
	service := siesta.NewService("")

	// Route accepts normal http.Handlers.
	service.Route("GET", "/", "Sends 'Hello, world!'", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	// Let's create some simple middleware.
	// This one will accept a Context argument.
	timestamper := func(c siesta.Context, w http.ResponseWriter, r *http.Request) {
		c.Set("start", time.Now())
	}

	timeHandler := func(c siesta.Context, w http.ResponseWriter, r *http.Request) {
		delta := time.Now().Sub(c.Get("start").(time.Time))
		fmt.Fprintf(w, "That took %v.\n", delta)
	}

	// We can compose the two handlers together and add it as a new route.
	service.Route("GET", "/time", "Sends how long it took to send a message", siesta.Compose(timestamper, timeHandler))

	// service is an http.Handler, so we can pass it directly to ListenAndServe.
	log.Fatal(http.ListenAndServe(":8080", service))
}
