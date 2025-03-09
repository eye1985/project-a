package middleware

import (
	"net/http"
)

func compose(fns ...Handler) Handler {
	if len(fns) == 0 {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return next
		}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(fns) - 1; i >= 0; i-- {
			next = fns[i](next)
		}

		return next
	}
}

type RouteHandler interface {
	HandleFunc(route string, h http.HandlerFunc)
	Handle(pattern string, handler http.Handler)
}

type Handler func(next http.HandlerFunc) http.HandlerFunc
type Middleware struct {
	middlewares    []Handler
	composedHandle Handler
	Mux            *http.ServeMux
}

func NewMiddlewareMux() *Middleware {
	return &Middleware{
		Mux:            http.NewServeMux(),
		middlewares:    []Handler{},
		composedHandle: compose(),
	}
}

func (m *Middleware) Add(mw Handler) {
	m.middlewares = append(m.middlewares, mw)
	m.composedHandle = compose(m.middlewares...)
}

func (m *Middleware) HandleFunc(route string, h http.HandlerFunc) {
	m.Mux.HandleFunc(route, m.composedHandle(h))
}

// Handle TODO add logging
func (m *Middleware) Handle(pattern string, handler http.Handler) {
	m.Mux.Handle(pattern, handler)
}
