package models

type Route struct {
	Id                 string `json:"id" example:"route1"`
	Method             string `json:"method" example:"GET"`
	Path               string `json:"path" example:"/users/:id"`
	QueryData          string `json:"queryData" example:"{page:1, limit:10}"`
	ParamData          string `json:"paramData" example:"{id:123}"`
	BodyData           string `json:"bodyData" example:"{name:John, age:30}"`
	ExpectedStatusCode int    `json:"expectedStatusCode" example:"200"`
	ExpectedBodyData   string `json:"expectedBodyData" example:"{status:success, message:User found}"`
}
