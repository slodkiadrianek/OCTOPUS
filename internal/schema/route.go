package schema

import (
	z "github.com/Oudwins/zog"
)

var RouteIDSchema = z.Struct(z.Shape{
	"routeID": z.String().Required(),
})

var CreateRouteSchema = z.Struct(z.Shape{
	"name": z.String().Required(),
	"routes": z.Slice(z.Struct(z.Shape{
		"path":                 z.String().Required().Max(256),
		"method":               z.String().OneOf([]string{"POST", "GET", "PUT", "PATCH", "DELETE"}),
		"requestAuthorization": z.String().Optional(),
		"requestQuery": z.CustomFunc[map[string]string](func(val *map[string]string, ctx z.Ctx) bool {
			return true
		}),
		"requestParams": z.CustomFunc[map[string]string](func(val *map[string]string, ctx z.Ctx) bool {
			return true
		}),
		"requestBody": z.CustomFunc[map[string]any](func(val *map[string]any, ctx z.Ctx) bool {
			return true
		}),
		"nextRouteBody":                z.Slice(z.String()).Optional(),
		"nextRouteQuery":               z.Slice(z.String()).Optional(),
		"nextRouteParams":              z.Slice(z.String()).Optional(),
		"nextRouteAuthorizationHeader": z.String().Optional(),
		"responseStatusCode":           z.Int().Required(),
		"responseBody": z.CustomFunc[map[string]any](func(val *map[string]any, ctx z.Ctx) bool {
			return true
		}),
	})),
})
