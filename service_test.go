package siesta

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServiceRoute(t *testing.T) {
	s := NewService("foos")
	s.Route(http.MethodGet, "/bars/:id/bazs", "Handles bars' bazs", func(Context, http.ResponseWriter, *http.Request, func()) {})

	srv := httptest.NewServer(s)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/foos/bars/1/bazs")
	if err != nil {
		t.Fatal(err)
	}
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("expected status %d got %d", want, got)
	}
}

func TestServiceDefaultNotFound(t *testing.T) {
	s := NewService("")

	srv := httptest.NewServer(s)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/no/where/to/be/found")
	if err != nil {
		t.Fatal(err)
	}
	if want, got := http.StatusNotFound, resp.StatusCode; want != got {
		t.Fatalf("expected status %d got %d", want, got)
	}
}

func TestServiceCustomNotFound(t *testing.T) {
	type payload struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	want := payload{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}

	s := NewService("")
	s.SetNotFound(func(c Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&want)
	})

	srv := httptest.NewServer(s)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/no/where/to/be/found")
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		t.Fatal(err)
	}
	if want, got := http.StatusNotFound, resp.StatusCode; want != got {
		t.Fatalf("expected status %d got %d", want, got)
	}

	var got payload
	if err = json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected payload %v got %v", want, got)
	}
}
