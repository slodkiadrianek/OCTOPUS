package schema

import z "github.com/Oudwins/zog"

var CreateRouteSchema = z.Slice(z.Struct(z.Shape{
	"path":               z.String().Required(),
	"method":             z.String().OneOf([]string{"POST", "GET", "PUT", "PATCH", "DELETE"}),
	"query_data":         z.String().Optional(),
	"param_data":         z.String().Optional(),
	"body_data":          z.String().Optional(),
	"expectedStatusCode": z.Int().Required(),
	"expectedBodyData":   z.String().Required(),
}))
