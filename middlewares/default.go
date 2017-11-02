package middlewares

import "net/http"

var DefaultMiddlewareStack = &MiddlewareStack{}

// Use use middleware with DefaultMiddlewareStack
func Use(middleware Middleware) {
	DefaultMiddlewareStack.Use(middleware)
}

// Remove remove middleware by name with DefaultMiddlewareStack
func Remove(name string) {
	DefaultMiddlewareStack.Remove(name)
}

// Apply apply DefaultMiddlewareStack's middlewares to handler
func Apply(handler http.Handler) http.Handler {
	return DefaultMiddlewareStack.Apply(handler)
}
