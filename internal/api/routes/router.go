package routes

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/middleware"
)

type (
	Router     struct{}
	Middleware func(http.Handler) http.Handler
)

func (r *Router) Request(route string, method string, fns ...any) {
	middlewares := []Middleware{}
	var finalHandler http.Handler
	for _, el := range fns {
		switch fn := el.(type) {
		case func(http.Handler) http.Handler:
			middlewares = append(middlewares, fn)
		case http.Handler:
			finalHandler = fn
		}
	}
	middlewares = append(middlewares, middleware.ErrorHandler)
	handler := finalHandler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	chainedHandler := middleware.CorsHandler(middleware.MethodCheck(handler, method))

	http.Handle(route, chainedHandler)
}

func (r *Router) Get(route string, fns ...any) {
	r.Request(route, "GET", fns)
}

func (r *Router) Post(route string, fns ...Middleware) {
	r.Request(route, "POST", fns)
}

func (r *Router) Patch(route string, fns ...Middleware) {
	r.Request(route, "PATCH", fns)
}

func (r *Router) PUT(route string, fns ...Middleware) {
	r.Request(route, "PUT", fns)
}

func (r *Router) Delete(route string, fns ...Middleware) {
	r.Request(route, "Delete", fns)
}
