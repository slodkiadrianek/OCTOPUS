package middleware

import (
	"fmt"
	"net/http"

	"github.com/slodkiadrianek/octopus/internal/Models"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

//	func ErrorHandler(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			errVal := r.Context().Value("Error")
//			err, ok := errVal.(error)
//			if ok && err != nil {
//				if customErr, isCustomErr := err.(*Models.Error); isCustomErr {
//					utils.SendResponse(w, customErr.StatusCode, map[string]string{"errorCategory": customErr.Category, "errorDescription": customErr.Descritpion})
//					return
//				}
//				utils.SendResponse(w, 500, map[string]string{"errorCategory": "Server", "errorDescription": "Internal server error"})
//				return
//			}
//			next.ServeHTTP(w, r)
//		})
//	}
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errVal := r.Context().Value("Error")
		if errVal == nil {
			if next != nil {

				fmt.Println("Calling next handler") // add this line
				next.ServeHTTP(w, r)
				fmt.Println("Next handler finished") // add this line
			}
			return
		}

		err, ok := errVal.(error)
		if !ok || err == nil {
			fmt.Println("Calling next handler") // add this line
			next.ServeHTTP(w, r)
			fmt.Println("Next handler finished") // add this line
			return
		}

		customErr, isCustomErr := err.(*Models.Error)
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
