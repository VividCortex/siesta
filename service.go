package siesta

import (
	"net/http"
	"path"
	"strings"
)

// Registered services keyed by base URI.
var services = map[string]*Service{}

// A Service is a container for routes with a common base URI.
// It also has two middleware chains, named "pre" and "post".
//
// The "pre" chain is run before the main handler. The first
// handler in the "pre" chain is guaranteed to run, but execution
// may quit anywhere else in the chain.
//
// If the "pre" chain executes completely, the main handler is executed.
// It is skipped otherwise.
//
// The "post" chain runs after the main handler, whether it is skipped
// or not. The first handler in the "post" chain is guaranteed to run, but
// execution may quit anywhere else in the chain if the quit function
// is called.
type Service struct {
	baseURI   string
	trimSlash bool

	pre  []contextHandler
	post []contextHandler

	routes map[string]*node

	notFound contextHandler
}

// NewService returns a new Service with the given base URI
// or panics if the base URI has already been registered.
func NewService(baseURI string) *Service {
	if services[baseURI] != nil {
		panic("service already registered")
	}

	return &Service{
		baseURI:   path.Join("/", baseURI, "/"),
		routes:    map[string]*node{},
		trimSlash: true,
	}
}

// DisableTrimSlash disables the removal of trailing slashes
// before route matching.
func (s *Service) DisableTrimSlash() {
	s.trimSlash = false
}

func addToChain(f interface{}, chain []contextHandler) []contextHandler {
	m := toContextHandler(f)
	return append(chain, m)
}

// AddPre adds f to the end of the "pre" chain.
// It panics if f cannot be converted to a contextHandler (see Service.Route).
func (s *Service) AddPre(f interface{}) {
	s.pre = addToChain(f, s.pre)
}

// AddPost adds f to the end of the "post" chain.
// It panics if f cannot be converted to a contextHandler (see Service.Route).
func (s *Service) AddPost(f interface{}) {
	s.post = addToChain(f, s.post)
}

// Service satisfies the http.Handler interface.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.ServeHTTPInContext(NewSiestaContext(), w, r)
}

// ServiceHTTPInContext serves an HTTP request within the Context c.
// A Service will run through both of its internal chains, quitting
// when requested.
func (s *Service) ServeHTTPInContext(c Context, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	quit := false
	for _, m := range s.pre {
		m(c, w, r, func() {
			quit = true
		})

		if quit {
			// Break out of the "pre" loop, but
			// continue on.
			break
		}
	}

	if !quit {
		// The main handler is only run if we have not
		// been signaled to quit.

		if r.URL.Path != "/" && s.trimSlash {
			r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		}

		var (
			handler contextHandler
			usage   string
			params  routeParams
		)

		// Lookup the tree for this method
		routeNode, ok := s.routes[r.Method]

		if ok {
			handler, usage, params, _ = routeNode.getValue(r.URL.Path)
			c.Set(UsageContextKey, usage)
		}

		if handler == nil {
			if s.notFound != nil {
				// Use user-defined handler.
				s.notFound(c, w, r, func() {})
			} else {
				// Default to the net/http NotFoundHandler.
				http.NotFoundHandler().ServeHTTP(w, r)
			}
		} else {
			for _, p := range params {
				r.Form.Set(p.Key, p.Value)
			}

			handler(c, w, r, func() {
				quit = true
			})
		}
	}

	quit = false
	for _, m := range s.post {
		m(c, w, r, func() {
			quit = true
		})

		if quit {
			return
		}
	}
}

// Route adds a new route to the Service.
// f must be a function with one of the following signatures:
//
//     func(http.ResponseWriter, *http.Request)
//     func(http.ResponseWriter, *http.Request, func())
//     func(Context, http.ResponseWriter, *http.Request)
//     func(Context, http.ResponseWriter, *http.Request, func())
//
// Note that Context is an interface type defined in this package.
// The last argument is a function which is called to signal the
// quitting of the current execution sequence.
func (s *Service) Route(verb, uriPath, usage string, f interface{}) {
	handler := toContextHandler(f)

	if n := s.routes[verb]; n == nil {
		s.routes[verb] = &node{}
	}

	s.routes[verb].addRoute(
		path.Join(s.baseURI, strings.TrimRight(uriPath, "/")),
		usage, handler)
}

// SetNotFound sets the handler for all paths that do not
// match any existing routes. It accepts the same function
// signatures that Route does with the addition of `nil`.
func (s *Service) SetNotFound(f interface{}) {
	if f == nil {
		s.notFound = nil
		return
	}

	handler := toContextHandler(f)
	s.notFound = handler
}

// Register registers s by adding it as a handler to the
// DefaultServeMux in the net/http package.
func (s *Service) Register() {
	http.Handle(s.baseURI, s)
}
