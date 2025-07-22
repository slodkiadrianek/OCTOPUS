package routes

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/utils"
)

type (
	Router     struct{}
	Middleware func(http.Handler) http.Handler
)

func methodCheck(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			utils.SendResponse(w, 404, map[string]string{"error": "Not found"})
		}
		next.ServeHTTP(w, r)
	})
}

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
	handler := finalHandler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	http.Handle(route, methodCheck(handler, method))
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
