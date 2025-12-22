package middleware

import (
	"bytes"
	"io"
	"net/http"
	"reflect"

	"github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/utils"
	"github.com/slodkiadrianek/octopus/internal/utils/request"
	"github.com/slodkiadrianek/octopus/internal/utils/response"
	"github.com/slodkiadrianek/octopus/internal/utils/validation"
)

func ValidateMiddleware[dataFromRequestType any, validationSchemaType *zog.StructSchema | *zog.SliceSchema](
	validationType string, validationSchema validationSchemaType,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return validateHandler[dataFromRequestType](next, validationType, validationSchema)
	}
}

func validateDataFromRequest[dataFromRequestType any, validationSchemaType *zog.StructSchema | *zog.SliceSchema](
	dataFromRequest *dataFromRequestType, validationSchema validationSchemaType,
) error {
	typeOfDataFromRequest := reflect.TypeOf(dataFromRequest)
	var errMap zog.ZogIssueMap
	if typeOfDataFromRequest.Kind() == reflect.Slice {
		// schema, _ := any(validationSchema).(*zog.SliceSchema)
		// errMap = validation.ValidateInputStruct(schema, dataFromRequest)
	} else {
		schema, _ := any(validationSchema).(*zog.StructSchema)
		errMap = validation.ValidateInputStruct(schema, dataFromRequest)
	}

	if errMap != nil {
		parsedError, err := utils.MarshalData(errMap["$first"])
		if err != nil {
			return err
		}
		return models.NewError(422, "Validation", string(parsedError))
	}

	return nil
}

func validateHandler[dataFromRequestType any, validationSchemaType *zog.StructSchema | *zog.SliceSchema](next http.Handler,
	validationType string,
	validationSchema validationSchemaType,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch validationType {
		case "body":
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				response.SetError(w, r, err)
				return
			}
			defer r.Body.Close()

			var dataFromRequest *dataFromRequestType
			dataFromRequest, err = utils.UnmarshalData[dataFromRequestType](bodyBytes)
			if err != nil {
				response.SetError(w, r, err)
				return
			}

			err = validateDataFromRequest[dataFromRequestType, validationSchemaType](dataFromRequest, validationSchema)
			if err != nil {
				response.SetError(w, r, err)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		case "params":
			paramsMap, err := request.ReadAllParams(r)
			if err != nil {
				response.SetError(w, r, err)
				return
			}

			paramBytes, err := utils.MarshalData(paramsMap)
			if err != nil {
				response.SetError(w, r, err)
				return
			}

			var dataFromRequest *dataFromRequestType
			dataFromRequest, err = utils.UnmarshalData[dataFromRequestType](paramBytes)
			if err != nil {
				response.SetError(w, r, err)
				return
			}
			err = validateDataFromRequest[dataFromRequestType, validationSchemaType](dataFromRequest, validationSchema)
			if err != nil {
				response.SetError(w, r, err)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
