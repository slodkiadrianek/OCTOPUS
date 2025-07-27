package routes

import (
	"fmt"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

type routeKey struct {
	method string
	path   string
}
type (
	Router struct {
		MiddlewarePreChain []Middleware
		routes             map[routeKey]http.Handler
	}
	Middleware func(http.Handler) http.Handler
)

func NewRouter() *Router {
	return &Router{
		MiddlewarePreChain: []Middleware{},
		routes:             make(map[routeKey]http.Handler),
	}
}

func (r *Router) Request(route string, method string, fns ...any) {
	middlewares := []Middleware{}
	middlewares = append(middlewares, middleware.MethodCheckMiddleware(method))
	var finalHandler http.Handler
	if len(r.MiddlewarePreChain) > 0 {
		middlewares = append(middlewares, r.MiddlewarePreChain...)
	}
	for _, el := range fns {
		switch fn := el.(type) {
		case func(http.Handler) http.Handler:
			middlewares = append(middlewares, fn)
		case func(http.ResponseWriter, *http.Request):
			finalHandler = http.HandlerFunc(fn)
		}
	}
	middlewares = append(middlewares, middleware.ErrorHandler)
	handler := finalHandler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	chainedHandler := handler
	r.routes[routeKey{method: method, path: route}] = chainedHandler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Incoming request: %s %s\n", req.Method, req.URL.Path)
	key := routeKey{method: req.Method, path: req.URL.Path}
	handler, ok := r.routes[key]
	fmt.Println(r.routes)
	if !ok {
		fmt.Printf("Route not found for key: %v\n", key)
		utils.SendResponse(w, 404, map[string]string{"errorDescription": "Route not found"})
		return
	}
	fmt.Println("Found matching handler, serving request")
	handler.ServeHTTP(w, req)
}

func (r *Router) Group(prefix string) *GroupRouter {
	return NewGroupRouter(prefix, r)
}

func (r *Router) USE(fns Middleware) {
	r.MiddlewarePreChain = append(r.MiddlewarePreChain, fns)
}

func (r *Router) GET(route string, fns ...any) {
	r.Request(route, http.MethodGet, fns...)
}

func (r *Router) POST(route string, fns ...any) {
	r.Request(route, http.MethodPost, fns...)
}

func (r *Router) PATCH(route string, fns ...any) {
	r.Request(route, http.MethodPatch, fns...)
}

func (r *Router) PUT(route string, fns ...any) {
	r.Request(route, http.MethodPut, fns...)
}

func (r *Router) DELETE(route string, fns ...any) {
	r.Request(route, http.MethodDelete, fns...)
}
