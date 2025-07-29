package middleware

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/utils"
)

func MethodCheckMiddleware(method string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return MethodCheckHandler(next, method)
	}
}

func MethodCheckHandler(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			utils.SendResponse(w, 405, map[string]string{"error": "Not found"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
