package siesta

import (
	"testing"
)

func TestContext(t *testing.T) {
	var c Context = NewSiestaContext()
	c.Set("foo", "bar")
	v := c.Get("foo")
	if v == nil {
		t.Fatal("expected to see a value for key `foo`")
	}

	if v.(string) != "bar" {
		t.Errorf("expected value %v, got %v", "bar", v.(string))
	}
}

func TestEmptyContext(t *testing.T) {
	var c Context = emptyContext{}
	c.Set("foo", "bar")
	v := c.Get("foo")
	if v != nil {
		t.Fatal("expected to not see a value for key `foo`")
	}
}
