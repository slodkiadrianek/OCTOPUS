package schema

import z "github.com/Oudwins/zog"

type CreateRoute struct {
	Path               string `json:"path" example:"/users"`
	Method             string `json:"method" example:"GET"`
	QueryData          string `json:"queryData" example:"id=1"`
	ParamData          string `json:"paramData" example:"id=1"`
	BodyData           string `json:"bodyData" example:"id=1"`
	ExpectedStatusCode int    `json:"expectedStatusCode" example:"200"`
	ExpectedBodyData   string `json:"expectedBodyData" example:"id=1"`
}

var CreateRouteSchema = z.Slice(z.Struct(z.Shape{
	"path":               z.String().Required(),
	"method":             z.String().OneOf([]string{"POST", "GET", "PUT", "PATCH", "DELETE"}),
	"query_data":         z.String().Optional(),
	"param_data":         z.String().Optional(),
	"body_data":          z.String().Optional(),
	"expectedStatusCode": z.Int().Required(),
	"expectedBodyData":   z.String().Required(),
}))
