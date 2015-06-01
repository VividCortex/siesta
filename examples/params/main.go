package main

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/VividCortex/siesta"
)

func main() {
	// Create a new Service rooted at "/"
	service := siesta.NewService("/")

	// Here's a handler that uses a URL parameter.
	// Example: GET /greet/Bob
	service.Route("GET", "/greet/:name", "Greets with a name.",
		func(w http.ResponseWriter, r *http.Request) {
			var params siesta.Params
			name := params.String("name", "", "Person's name")

			err := params.Parse(r.Form)
			if err != nil {
				log.Println("Error parsing parameters!", err)
				return
			}

			fmt.Fprintf(w, "Hello, %s!", *name)
		},
	)

	// Here's a handler that uses a query string parameter.
	// Example: GET /square?number=10
	service.Route("GET", "/square", "Prints the square of a number.",
		func(w http.ResponseWriter, r *http.Request) {
			var params siesta.Params
			number := params.Int("number", 0, "A number to square")

			err := params.Parse(r.Form)
			if err != nil {
				log.Println("Error parsing parameters!", err)
				return
			}

			fmt.Fprintf(w, "%d * %d = %d.", *number, *number, (*number)*(*number))
		},
	)

	// We can also use both URL and query string parameters.
	// Example: GET /exponentiate/10?power=10
	service.Route("GET", "/exponentiate/:number", "Exponentiates a number.",
		func(w http.ResponseWriter, r *http.Request) {
			var params siesta.Params
			number := params.Float64("number", 0, "A number to exponentiate")
			power := params.Float64("power", 1, "Power")

			err := params.Parse(r.Form)
			if err != nil {
				log.Println("Error parsing parameters!", err)
				return
			}

			fmt.Fprintf(w, "%g ^ %g = %g.", *number, *power, math.Pow(*number, *power))
		},
	)

	// service is an http.Handler, so we can pass it directly to ListenAndServe.
	log.Fatal(http.ListenAndServe(":8080", service))
}
