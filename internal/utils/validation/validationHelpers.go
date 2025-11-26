package validation

import (
	z "github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
)

func ValidateUsersIds(userId, userIdToken int) error {
	if userIdToken != userId {
		err := models.NewError(500, "Server", "Internal server error")
		return err
	}
	return nil
}

func ValidateInputStruct(schema *z.StructSchema, val any) z.ZogIssueMap {
	errMap := schema.Validate(val)
	if errMap != nil {
		return errMap
	}
	return nil
}

func ValidateInputSlice(schema *z.SliceSchema, val any) z.ZogIssueMap {
	errMap := schema.Validate(val)
	if errMap != nil {
		return errMap
	}
	return nil
}

func CheckIsNextRouteBodyInTheBodyAndInTheBodyOfTheNextRoute(actualRoute DTO.CreateRoute, nextRoute DTO.CreateRoute) bool {
	for _, val := range actualRoute.NextRouteBody {
		resBody := IsDataInResponseOrRequest[map[string]any](actualRoute.ResponseBody, val)
		if !resBody {
			return false
		}
		reqBody := IsDataInResponseOrRequest[map[string]any](nextRoute.RequestBody, val)
		if !reqBody {
			return false
		}
	}
	return true
}

func CheckIsNextRouteQueryInTheBodyAndInTheQueryOfTheNextRoute(actualRoute DTO.CreateRoute, nextRoute DTO.CreateRoute) bool {
	for _, val := range actualRoute.NextRouteQuery {
		resBody := IsDataInResponseOrRequest[map[string]any](actualRoute.ResponseBody, val)
		if !resBody {
			return false
		}
		reqQuery := IsDataInResponseOrRequest[map[string]string](nextRoute.RequestQuery, val)
		if !reqQuery {
			return false
		}
	}
	return true
}

func CheckIsNextRouteParamsInTheBodyAndInTheParamsOfTheNextRoute(actualRoute DTO.CreateRoute, nextRoute DTO.CreateRoute) bool {
	for _, val := range actualRoute.NextRouteParams {
		resBody := IsDataInResponseOrRequest[map[string]any](actualRoute.ResponseBody, val)
		if !resBody {
			return false
		}
		reqParams := IsDataInResponseOrRequest[map[string]string](nextRoute.RequestParams, val)
		if !reqParams {
			return false
		}
	}
	return true
}

func IsDataInResponseOrRequest[T map[string]any | map[string]string](responseBody T, data string) bool {
	switch body := any(responseBody).(type) {
	case map[string]any:
		for key, val := range body {
			if key == data {
				return true
			}
			if x, ok := val.(map[string]any); ok {
				if IsDataInResponseOrRequest(x, data) {
					return true
				}
			}
		}
	case map[string]string:
		for key := range body {
			if key == data {
				return true
			}
		}
	}
	return false
}
