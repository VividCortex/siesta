package siesta

import (
	"errors"
	"net/http"
)

var ErrUnsupportedHandler = errors.New("siesta: unsupported handler")

type contextHandler func(Context, http.ResponseWriter, *http.Request, func())

func (h contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(emptyContext{}, w, r, nil)
}

func (h contextHandler) ServeHTTPInContext(c Context, w http.ResponseWriter, r *http.Request) {
	h(c, w, r, nil)
}

func toContextHandler(f interface{}) contextHandler {
	var m contextHandler

	switch f.(type) {
	case func(Context, http.ResponseWriter, *http.Request, func()):
		m = contextHandler(f.(func(Context, http.ResponseWriter, *http.Request, func())))
	case contextHandler:
		m = f.(contextHandler)
	case func(Context, http.ResponseWriter, *http.Request):
		m = func(c Context, w http.ResponseWriter, r *http.Request, q func()) {
			f.(func(Context, http.ResponseWriter, *http.Request))(c, w, r)
		}
	case func(http.ResponseWriter, *http.Request, func()):
		m = func(c Context, w http.ResponseWriter, r *http.Request, q func()) {
			f.(func(http.ResponseWriter, *http.Request, func()))(w, r, q)
		}
	case func(http.ResponseWriter, *http.Request):
		m = func(c Context, w http.ResponseWriter, r *http.Request, q func()) {
			f.(func(http.ResponseWriter, *http.Request))(w, r)
		}
	default:

		// Check for http.Handlers too.
		if h, ok := f.(http.Handler); ok {
			return toContextHandler(h.ServeHTTP)
		}

		panic(ErrUnsupportedHandler)
	}

	return m
}

// Compose composes multiple contextHandlers into a single contextHandler.
func Compose(stack ...interface{}) contextHandler {
	contextStack := make([]contextHandler, 0, len(stack))
	for i := range stack {
		m := toContextHandler(stack[i])

		contextStack = append(contextStack, m)
	}

	return func(c Context, w http.ResponseWriter, r *http.Request, quit func()) {
		quitStack := false

		for _, m := range contextStack {
			m(c, w, r, func() {
				quitStack = true
			})

			if quitStack {
				quit()
				break
			}
		}
	}
}
