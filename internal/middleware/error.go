package middleware

import (
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/Models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errVal := r.Context().Value("Error")
		err, ok := errVal.(error)
		if ok && err != nil {
			if customErr, isCustomErr := err.(*Models.Error); isCustomErr {
				utils.SendResponse(w, customErr.StatusCode, map[string]string{"errorCategory": customErr.Category, "errorDescription": customErr.Descritpion})
				return
			}
			utils.SendResponse(w, 500, map[string]string{"errorCategory": "Server", "errorDescription": "Internal server error"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
