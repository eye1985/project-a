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

type Handler func(next http.HandlerFunc) http.HandlerFunc
type Middleware struct {
	middlewares    []Handler
	composedHandle Handler
	Mux            *http.ServeMux
}

func NewMux() *Middleware {
	return &Middleware{
		Mux:            http.NewServeMux(),
		middlewares:    []Handler{},
		composedHandle: compose(),
	}
}

// Add Inserts Global middlewares, global middlewares are run after the handler middlewares
func (m *Middleware) Add(mw Handler) {
	m.middlewares = append(m.middlewares, mw)
	m.composedHandle = compose(m.middlewares...)
}

func (m *Middleware) HandleFunc(route string, h http.HandlerFunc, middlewareHandlers ...Handler) {
	middlewares := append(middlewareHandlers, m.middlewares...)
	composed := compose(middlewares...)
	m.Mux.HandleFunc(route, composed(h))
}

// Handle TODO add logging
func (m *Middleware) Handle(pattern string, handler http.Handler) {
	m.Mux.Handle(pattern, handler)
}
