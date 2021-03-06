package siesta

import (
	"errors"
	"net/http"
)

var ErrUnsupportedHandler = errors.New("siesta: unsupported handler")

// ContextHandler is a siesta handler.
type ContextHandler func(Context, http.ResponseWriter, *http.Request, func())

func (h ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(EmptyContext{}, w, r, nil)
}

func (h ContextHandler) ServeHTTPInContext(c Context, w http.ResponseWriter, r *http.Request) {
	h(c, w, r, nil)
}

// ToContextHandler transforms f into a ContextHandler.
// f must be a function with one of the following signatures:
//     func(http.ResponseWriter, *http.Request)
//     func(http.ResponseWriter, *http.Request, func())
//     func(Context, http.ResponseWriter, *http.Request)
//     func(Context, http.ResponseWriter, *http.Request, func())
func ToContextHandler(f interface{}) ContextHandler {
	switch t := f.(type) {
	case func(Context, http.ResponseWriter, *http.Request, func()):
		return ContextHandler(t)
	case ContextHandler:
		return t
	case func(Context, http.ResponseWriter, *http.Request):
		return func(c Context, w http.ResponseWriter, r *http.Request, q func()) {
			t(c, w, r)
		}
	case func(http.ResponseWriter, *http.Request, func()):
		return func(c Context, w http.ResponseWriter, r *http.Request, q func()) {
			t(w, r, q)
		}
	case func(http.ResponseWriter, *http.Request):
		return func(c Context, w http.ResponseWriter, r *http.Request, q func()) {
			t(w, r)
		}
	case http.Handler:
		return func(c Context, w http.ResponseWriter, r *http.Request, q func()) {
			t.ServeHTTP(w, r)
		}
	default:
		panic(ErrUnsupportedHandler)
	}
}

// Compose composes multiple ContextHandlers into a single ContextHandler.
func Compose(stack ...interface{}) ContextHandler {
	contextStack := make([]ContextHandler, 0, len(stack))
	for i := range stack {
		m := ToContextHandler(stack[i])

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
