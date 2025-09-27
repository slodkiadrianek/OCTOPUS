package DTO

type CreateRoute struct {
	Path               string `json:"path" example:"/users"`
	Method             string `json:"method" example:"GET"`
	QueryData          string `json:"queryData" example:"id=1"`
	ParamData          string `json:"paramData" example:"id=1"`
	BodyData           string `json:"bodyData" example:"id=1"`
	ExpectedStatusCode int    `json:"expectedStatusCode" example:"200"`
	ExpectedBodyData   string `json:"expectedBodyData" example:"id=1"`
}

func NewCreateRoute(method, path, queryData, paramData, bodyData string, expectedStatusCode int, predictedBodyData string, appId int) *CreateRoute {
	return &CreateRoute{
		Method:             method,
		Path:               path,
		QueryData:          queryData,
		ParamData:          paramData,
		BodyData:           bodyData,
		ExpectedStatusCode: expectedStatusCode,
		ExpectedBodyData:   predictedBodyData,
	}
}
