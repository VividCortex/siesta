package siesta

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompose(t *testing.T) {
	key := "i"
	stack := Compose(
		func(c Context, w http.ResponseWriter, r *http.Request, quit func()) {
			r.Header.Set(key, r.Header.Get(key)+"a")
			i, _ := c.Get(key).(int)
			c.Set(key, i+2)
		},
		func(c Context, w http.ResponseWriter, r *http.Request) {
			r.Header.Set(key, r.Header.Get(key)+"b")
			i, _ := c.Get(key).(int)
			c.Set(key, i+4)
		},
		func(w http.ResponseWriter, r *http.Request, quit func()) {
			r.Header.Set(key, r.Header.Get(key)+"c")
		},
		func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(key, r.Header.Get(key)+"d")
		},
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(key, r.Header.Get(key)+"e")
		}),
	)

	c := NewSiestaContext()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	stack.ServeHTTPInContext(c, w, r)

	i, _ := c.Get(key).(int)
	if want, got := 6, i; want != got {
		t.Errorf("expected %d got %d", want, got)
	}
	if want, got := "abcde", r.Header.Get(key); want != got {
		t.Errorf("expected %s got %s", want, got)
	}
}

func TestToContextHandlerUnsupportedHandler(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected a panic")
		}
		err, _ := r.(error)
		if want, got := ErrUnsupportedHandler, err; want != got {
			t.Fatalf("expected %v got %v", want, got)
		}
	}()

	_ = ToContextHandler(func() {})
}
