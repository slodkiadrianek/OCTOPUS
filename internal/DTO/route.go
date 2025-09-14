package DTO

type CreateRoute struct {
	Id                 string `json:"id" example:"route1"`
	Method             string `json:"method" example:"GET"`
	Path               string `json:"path" example:"/users/:id"`
	QueryData          string `json:"queryData" example:"{page:1, limit:10}"`
	ParamData          string `json:"paramData" example:"{id:123}"`
	BodyData           string `json:"bodyData" example:"{name:John, age:30}"`
	ExpectedStatusCode int    `json:expectedStatusCode" example:"200"`
	ExpectedBodyData   string `json:expectedBodyData" example:"{name:John, age:30}"`
}

func NewCreateRoute(id, method, path, queryData, paramData, bodyData string, expectedStatusCode int, predictedBodyData string, appId int) *CreateRoute {
	return &CreateRoute{
		Id:                 id,
		Method:             method,
		Path:               path,
		QueryData:          queryData,
		ParamData:          paramData,
		BodyData:           bodyData,
		ExpectedStatusCode: expectedStatusCode,
		ExpectedBodyData:   predictedBodyData,
	}
}
