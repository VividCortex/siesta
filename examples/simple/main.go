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
	service := siesta.NewService("/")

	// Route accepts normal http.Handlers.
	service.Route("GET", "/", "Sends 'Hello, world!'", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	// Let's create some simple "middleware."
	// This handler will accept a Context argument and will add the current
	// time to it.
	timestamper := func(c siesta.Context, w http.ResponseWriter, r *http.Request) {
		c.Set("start", time.Now())
	}

	// This is the handler that will actually send data back to the client.
	// It also takes a Context argument so it can get the timestamp from the
	// previous handler.
	timeHandler := func(c siesta.Context, w http.ResponseWriter, r *http.Request) {
		start := c.Get("start").(time.Time)
		delta := time.Now().Sub(start)
		fmt.Fprintf(w, "That took %v.\n", delta)
	}

	// We can compose these handlers together.
	timeHandlers := siesta.Compose(timestamper, timeHandler)

	// Finally, we'll add the new handler we created using composition to a new route.
	service.Route("GET", "/time", "Sends how long it took to send a message", timeHandlers)

	// service is an http.Handler, so we can pass it directly to ListenAndServe.
	log.Fatal(http.ListenAndServe(":8080", service))
}
