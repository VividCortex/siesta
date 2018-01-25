package siesta

// prepending nullByteStr avoids accidental key collisions
const nullByteStr = "\x00"

// UsageContextKey is a special context key to get the route usage information
// within a handler.
const UsageContextKey = nullByteStr + "usage"

// A siesta Context is a context interface that gets passed to each
// contextHandler.
type Context interface {
	Set(string, interface{})
	Get(string) interface{}
}

// This is a blank context.
type emptyContext struct{}

func (c emptyContext) Set(key string, value interface{}) {
}

func (c emptyContext) Get(key string) interface{} {
	return nil
}

// SiestaContext is a concrete implementation of the siesta.Context
// interface. Typically this will be created by the siesta framework
// itself upon each request. However creating your own SiestaContext
// might be useful for testing to isolate the behavior of a single
// handler.
type SiestaContext map[string]interface{}

func NewSiestaContext() SiestaContext {
	return SiestaContext{}
}

func (c SiestaContext) Set(key string, value interface{}) {
	c[key] = value
}

func (c SiestaContext) Get(key string) interface{} {
	return c[key]
}
