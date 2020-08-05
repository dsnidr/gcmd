package gcmd

// Command is a struct which holds the data required for a command
type Command struct {
	Name       string
	Usage      string
	middleware []MiddlewareFunc
	Handler    HandlerFunc
}

// Context holds context data for use by command handlers
type Context struct {
	Args    []string
	Store   Store
	Command *Command
}

// HandlerFunc is a type representation of a command handling function. Handler functions should return true if they
// were executed successfully, and false if they were cancelled.
type HandlerFunc func(c Context) error

// MiddlewareFunc represents a chainable middleware function
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Use adds middleware to the middleware chain
func (cmd *Command) Use(middleware ...MiddlewareFunc) {
	if cmd.middleware == nil {
		cmd.middleware = []MiddlewareFunc{}
	}

	cmd.middleware = append(cmd.middleware, middleware...)
}

func (cmd *Command) applyMiddleware() HandlerFunc {
	handler := cmd.Handler

	for i := len(cmd.middleware) - 1; i >= 0; i-- {
		handler = cmd.middleware[i](handler)
	}

	return handler
}

// Set retrieves an interface from the context
func (c *Context) Set(key string, value interface{}) {
	c.Store[key] = value
}

// Get retrieves an interface from the context
func (c *Context) Get(key string) interface{} {
	return c.Store[key]
}
