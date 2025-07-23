package middleware

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/utils"
)

func MethodCheck(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			utils.SendResponse(w, 404, map[string]string{"error": "Not found"})
		}
		next.ServeHTTP(w, r)
	})
}
