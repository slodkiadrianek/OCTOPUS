package middleware

import (
	"github.com/slodkiadrianek/octopus/internal/models"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/utils"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errVal := r.Context().Value("Error")
		if errVal == nil {
			if next != nil {
				next.ServeHTTP(w, r)
			}
			return
		}

		err, ok := errVal.(error)
		if !ok || err == nil {
			next.ServeHTTP(w, r)
			return
		}

		customErr, isCustomErr := err.(*models.Error)
		if isCustomErr {
			if customErr == nil {
				// Log or zwróć bezpieczną odpowiedź, aby uniknąć panic
				utils.SendResponse(w, 500, map[string]string{
					"errorCategory":    "Server",
					"errorDescription": "Internal server error (nil custom error)",
				})
				return
			}
			utils.SendResponse(w, customErr.StatusCode, map[string]string{
				"errorCategory":    customErr.Category,
				"errorDescription": customErr.Descritpion,
			})
			return
		}

		// Jeśli nie jest custom error, zwracamy 500
		utils.SendResponse(w, 500, map[string]string{
			"errorCategory":    "Server",
			"errorDescription": "Internal server error",
		})
	})
}
