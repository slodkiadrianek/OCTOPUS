package DTO

type RoutesParentId interface {
	GetParentId() int
	SetParentId(id int)
}
type CreateRoute struct {
	Path                 string            `json:"path" example:"/users"`
	Method               string            `json:"method" example:"GET"`
	RequestAuthorization string            `json:"requestAuthorization" example:"Bearer:fb43fg3487f34g78f3gu"`
	RequestQuery         map[string]string `json:"requestQuery" example:"id=1"`
	RequestParams        map[string]string `json:"requestParams" example:"id=1"`
	RequestBody          map[string]any    `json:"requestBody" example:"id=1"`
	NextRouteBody        []string          `json:"nextRouteBody"`
	NextRouteQuery       []string          `json:"nextRouteQuery"`
	NextRouteParams      []string          `json:"nextRouteParams"`
	ResponseStatusCode   int               `json:"responseStatusCode" example:"200"`
	ResponseBody         map[string]any    `json:"responseBody" example:"id=1"`
	ParentId             int
}

type WorkingRoute struct {
	ParentId        int
	AppId           string
	RouteId         int
	RequestId       int
	ResponseId      int
	NextRouteDataId int
	Status          string
}

type RouteInfo struct {
	Path     string `json:"path" example:"/users"`
	Method   string `json:"method" example:"GET"`
	ParentId int    `json:"parentId"`
}

func NewRouteInfo(path, method string, parentId int) *RouteInfo {
	return &RouteInfo{
		Path:     path,
		Method:   method,
		ParentId: parentId,
	}
}

func (ri *RouteInfo) GetParentId() int {
	return ri.ParentId
}

func (ri *RouteInfo) SetParentId(id int) {
	ri.ParentId = id
}

type RouteRequest struct {
	RequestAuthorization string `json:"requestAuthorization" example:"Bearer:fb43fg3487f34g78f3gu"`
	RequestQuery         string `json:"requestQuery" example:"id=1"`
	RequestParams        string `json:"requestParams" example:"id=1"`
	RequestBody          string `json:"requestBody" example:"id=1"`
	ParentId             int    `json:"parentId"`
}

func NewRouteRequest(requestAuthorization, requestQuery, requestParams, requestBody string, parentId int) *RouteRequest {
	return &RouteRequest{
		RequestAuthorization: requestAuthorization,
		RequestQuery:         requestQuery,
		RequestParams:        requestParams,
		RequestBody:          requestBody,
		ParentId:             parentId,
	}
}

func (rr *RouteRequest) GetParentId() int {
	return rr.ParentId
}

func (rr *RouteRequest) SetParentId(id int) {
	rr.ParentId = id
}

type NextRouteData struct {
	NextRouteBody   string `json:"nextRouteBody"`
	NextRouteQuery  string `json:"nextRouteQuery"`
	NextRouteParams string `json:"nextRouteParams"`
	ParentId        int    `json:"parentId"`
}

func NewNextRouteData(nextRouteBody, nextRouteQuery, nextRouteParams string, parentId int) *NextRouteData {
	return &NextRouteData{
		NextRouteBody:   nextRouteBody,
		NextRouteQuery:  nextRouteQuery,
		NextRouteParams: nextRouteParams,
		ParentId:        parentId,
	}
}

func (nrd *NextRouteData) GetParentId() int {
	return nrd.ParentId
}

func (nrd *NextRouteData) SetParentId(id int) {
	nrd.ParentId = id
}

type RouteResponse struct {
	ResponseStatusCode int    `json:"responseStatusCode" example:"200"`
	ResponseBody       string `json:"responseBody" example:"id=1"`
	ParentId           int    `json:"parentId"`
}

func NewRouteResponse(responseStatusCode, parentId int, responseBody string) *RouteResponse {
	return &RouteResponse{
		ResponseStatusCode: responseStatusCode,
		ResponseBody:       responseBody,
		ParentId:           parentId,
	}
}

func (rr *RouteResponse) GetParentId() int {
	return rr.ParentId
}

func (rr *RouteResponse) SetParentId(id int) {
	rr.ParentId = id
}
