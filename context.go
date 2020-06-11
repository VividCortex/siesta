package siesta

// Prepending nullByteStr avoids accidental context key collisions.
const nullByteStr = "\x00"

// UsageContextKey is a special context key to get the route usage information
// within a handler.
const UsageContextKey = nullByteStr + "usage"

// Context is a context interface that gets passed to each ContextHandler.
type Context interface {
	Set(string, interface{})
	Get(string) interface{}
}

// EmptyContext is a blank context.
type EmptyContext struct{}

func (c EmptyContext) Set(key string, value interface{}) {
}

func (c EmptyContext) Get(key string) interface{} {
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
