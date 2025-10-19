package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/utils"
)

func ValidateMiddleware[T any, Y *zog.StructSchema | *zog.SliceSchema](validationType string, validationSchema Y) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return validateHandler[T](next, validationType, validationSchema)
	}
}

func validateHandler[T any, Y *zog.StructSchema | *zog.SliceSchema](next http.Handler, validationType string, validationSchema Y) http.Handler {
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
			fmt.Println(err)
			if err != nil {
				utils.SendResponse(w, http.StatusBadRequest, "Invalid request body")
				return
			}
			var t T
			typeOfT := reflect.TypeOf(t)
			var errMap zog.ZogIssueMap
			if typeOfT.Kind() == reflect.Slice {
				schema, _ := any(validationSchema).(*zog.SliceSchema)
				errMap = utils.ValidateInputSlice(schema, data)
			} else {
				schema, _ := any(validationSchema).(*zog.StructSchema)
				errMap = utils.ValidateInputStruct(schema, data)
			}
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
			var t T
			typeOfT := reflect.TypeOf(t)
			var errMap zog.ZogIssueMap
			if typeOfT.Kind() == reflect.Slice {
				schema, _ := any(validationSchema).(*zog.SliceSchema)
				errMap = utils.ValidateInputSlice(schema, param)
			} else {
				schema, _ := any(validationSchema).(*zog.StructSchema)
				errMap = utils.ValidateInputStruct(schema, param)
			}
			if errMap != nil {
				utils.SendResponse(w, 422, errMap["$first"])
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
