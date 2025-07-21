package routes

import "net/http"

type (
	Router     struct{}
	Middleware func(http.Handler) http.Handler
)

func (r *Router) Get(route string, finalHandler http.Handler, fns ...Middleware) {
	handler := finalHandler
	for i := len(fns) - 1; i >= 0; i-- {
		handler = fns[i](handler)
	}
	http.HandleFunc(route, handler)
}

func (r *Router) Post(route string, fns ...Middleware) {
	for _, fn := range fns {
		http.HandleFunc(route, fn)
	}
}
