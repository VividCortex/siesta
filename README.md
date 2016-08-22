siesta [![Circle CI](https://circleci.com/gh/VividCortex/siesta.png?style=badge&circle-token=6b783c688fd8c3faed3554ca1e3548168ed87b10)](https://circleci.com/gh/VividCortex/siesta) [![GoDoc](https://godoc.org/github.com/VividCortex/siesta?status.svg)](https://godoc.org/github.com/VividCortex/siesta)
====

Siesta is a framework for writing composable HTTP handlers in Go. It supports typed URL parameters, middleware chains, and context passing.

Getting started
---
Siesta offers a `Service` type, which is a collection of middleware chains and handlers rooted at a base URI. There is no distinction between a middleware function and a handler function; they are all considered to be handlers and have access to the same arguments.

Siesta accepts many types of handlers. Refer to the [GoDoc](https://godoc.org/github.com/VividCortex/siesta#Service.Route) documentation for `Service.Route` for more information.

Here is the `simple` program in the examples directory. It demonstrates the use of a `Service`, routing, middleware, and a `Context`.

```go
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
	// The arguments are the method, path, description,
	// and the handler.
	service.Route("GET", "/", "Sends 'Hello, world!'",
		func(w http.ResponseWriter, r *http.Request) {
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
```

Siesta also provides utilities to manage URL parameters similar to the flag package. Refer to the `params` [example](https://github.com/VividCortex/siesta/blob/master/examples/params/main.go) for a demonstration.

Contributing
---
We only accept pull requests for minor fixes or improvements. This includes:

* Small bug fixes
* Typos
* Documentation or comments

Please open issues to discuss new features. Pull requests for new features will be rejected,
so we recommend forking the repository and making changes in your fork for your use case.

License
---
Siesta is licensed under the MIT license. The router, which is adapted from [httprouter](https://github.com/julienschmidt/httprouter), is licensed [separately](https://github.com/VividCortex/siesta/blob/6ce42bf31875cc845310b1f4775129edfc8d9967/tree.go#L2-L24).
