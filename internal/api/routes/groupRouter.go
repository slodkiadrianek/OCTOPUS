package routes

import (
	"net/http"
)

type GroupRouter struct {
	Prefix string
	router *Router
}

func NewGroupRouter(prefix string, router *Router) *GroupRouter {
	return &GroupRouter{
		Prefix: prefix,
		router: router,
	}
}

// func (rg *GroupRouter) USE(fns routes.Middleware) {
// 	rg.MiddlewarePreChain = append(rg.router.MiddlewarePreChain, fns)
// }

func (rg *GroupRouter) GET(route string, fns ...any) {
	rg.router.Request(rg.Prefix+route, http.MethodGet, fns...)
}

func (rg *GroupRouter) POST(route string, fns ...any) {
	rg.router.Request(rg.Prefix+route, http.MethodPost, fns...)
}

func (rg *GroupRouter) PATCH(route string, fns ...any) {
	rg.router.Request(rg.Prefix+route, http.MethodPatch, fns...)
}

func (rg *GroupRouter) PUT(route string, fns ...any) {
	rg.router.Request(rg.Prefix+route, http.MethodPut, fns...)
}

func (rg *GroupRouter) DELETE(route string, fns ...any) {
	rg.router.Request(rg.Prefix+route, http.MethodDelete, fns...)
}
