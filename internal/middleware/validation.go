package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

func ValidateMiddleware[T any](validationType string, validationSchema *zog.StructSchema) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return validateHandler[T](next, validationType, validationSchema)
	}
}

func validateHandler[T any](next http.Handler, validationType string, validationSchema *zog.StructSchema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch validationType {
		case "body":
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				utils.SendResponse(w, http.StatusBadRequest, "Invalid request body")
			}
			defer r.Body.Close()
			var data *T
			data, err = utils.UnmarshalData[T](bodyBytes)
			if err != nil {
				utils.SendResponse(w, http.StatusBadRequest, "Invalid request body")
				return
			}
			errMap := utils.ValidateInput(validationSchema, data)
			if errMap != nil {
				utils.SendResponse(w, 422, errMap["$first"])
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		case "params":
			paramsMap, err := utils.ReadAllParams(r)
			if err != nil {
				utils.SendResponse(w, http.StatusBadRequest, "Invalid request body")
				return
			}
			paramBytes, err := utils.MarshalData(paramsMap)
			if err != nil {
				utils.SendResponse(w, http.StatusBadRequest, "Invalid request parameters")
				return
			}

			param, err := utils.UnmarshalData[T](paramBytes)
			if err != nil {
				utils.SendResponse(w, http.StatusBadRequest, err)
				return
			}
			fmt.Println(param)
			errMap := utils.ValidateInput(validationSchema, param)
			if errMap != nil {
				utils.SendResponse(w, 422, errMap["$first"])
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
