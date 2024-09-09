package simplerest

// Middleware is a function that wraps a handler.
type Middleware func(Handler) Handler
