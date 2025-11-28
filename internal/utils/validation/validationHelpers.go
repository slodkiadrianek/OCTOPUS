package validation

import (
	"errors"
	z "github.com/Oudwins/zog"
	"github.com/slodkiadrianek/octopus/internal/DTO"
)

func ValidateUsersIds(userId, userIdFromToken int) error {
	if userIdFromToken != userId {
		err := errors.New("Provided user id's are different")
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
