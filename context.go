package siesta

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

type napContext struct {
	values map[string]interface{}
}

func newNapContext() *napContext {
	return &napContext{
		values: make(map[string]interface{}),
	}
}

func (c *napContext) Set(key string, value interface{}) {
	c.values[key] = value
}

func (c *napContext) Get(key string) interface{} {
	return c.values[key]
}
