package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errBucket := &models.ErrorBucket{}
		ctx := context.WithValue(r.Context(), "ErrorBucket", errBucket)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		errVal := errBucket.Err
		if errVal == nil {
			return
		}
		err, ok := errVal.(error)
		if !ok || err == nil {
			return
		}
		var customErr *models.Error
		isCustomErr := errors.As(err, &customErr)
		if isCustomErr {
			if customErr == nil {
				response.Send(w, 500, map[string]string{
					"errorCategory":    "Server",
					"errorDescription": "Internal server error",
				})
				return
			}
			response.Send(w, customErr.StatusCode, map[string]string{
				"errorCategory":    customErr.Category,
				"errorDescription": customErr.Description,
			})
			return
		}

		response.Send(w, 500, map[string]string{
			"errorCategory":    "Server",
			"errorDescription": "Internal server error",
		})
	})
}
