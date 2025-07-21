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

func (r *Router) Get(route string, fns ...interface{}) {
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
	http.Handle(route, methodCheck(handler, "GET"))
}

func (r *Router) Post(route string, fns ...Middleware) {
	handler := finalHandler
	for i := len(fns) - 1; i >= 0; i-- {
		handler = fns[i](handler)
	}
	http.Handle(route, handler)
}
