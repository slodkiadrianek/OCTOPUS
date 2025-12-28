package DTO

type RouteID struct {
	RouteID string `json:"routeID"`
}
type RoutesParentID interface {
	GetParentID() int
}
type CreateRouteData struct {
	Name   string `json:"name"`
	Routes []CreateRoute
}
type CreateRoute struct {
	Path                         string            `json:"path" example:"/users"`
	Method                       string            `json:"method" example:"GET"`
	RequestAuthorization         string            `json:"requestAuthorization" example:"Bearer:fb43fg3487f34g78f3gu"`
	RequestQuery                 map[string]string `json:"requestQuery" example:"id=1"`
	RequestParams                map[string]string `json:"requestParams" example:"id=1"`
	RequestBody                  map[string]any    `json:"requestBody" example:"id=1"`
	NextRouteBody                []string          `json:"nextRouteBody"`
	NextRouteQuery               []string          `json:"nextRouteQuery"`
	NextRouteParams              []string          `json:"nextRouteParams"`
	NextRouteAuthorizationHeader string            `json:"nextRouteAuthorizationHeader"`
	ResponseStatusCode           int               `json:"responseStatusCode" example:"200"`
	ResponseBody                 map[string]any    `json:"responseBody" example:"id=1"`
	ParentID                     int
}

type WorkingRoute struct {
	Name            string
	ParentID        int
	AppID           string
	RouteID         int
	RequestID       int
	ResponseID      int
	NextRouteDataID int
	Status          string
}

type RouteInfo struct {
	Path     string `json:"path" example:"/users"`
	Method   string `json:"method" example:"GET"`
	ParentID int    `json:"parentID"`
}

func NewRouteInfo(path, method string, parentID int) *RouteInfo {
	return &RouteInfo{
		Path:     path,
		Method:   method,
		ParentID: parentID,
	}
}

func (ri *RouteInfo) GetParentID() int {
	return ri.ParentID
}

type RouteRequest struct {
	AuthorizationHeader string `json:"authorizationHeader" example:"Bearer:fb43fg3487f34g78f3gu"`
	Query               string `json:"query" example:"id=1"`
	Params              string `json:"params" example:"id=1"`
	Body                string `json:"body" example:"id=1"`
	ParentID            int    `json:"parentID"`
}

func NewRouteRequest(authorizationHeader, query, params, body string, parentID int) *RouteRequest {
	return &RouteRequest{
		AuthorizationHeader: authorizationHeader,
		Query:               query,
		Params:              params,
		Body:                body,
		ParentID:            parentID,
	}
}

func (rr *RouteRequest) GetParentID() int {
	return rr.ParentID
}

type NextRoute struct {
	Body                string `json:"body"`
	Query               string `json:"query"`
	Params              string `json:"params"`
	AuthorizationHeader string `json:"authorizationHeader"`
	ParentID            int    `json:"parentID"`
}

func NewNextRouteData(body, query, params, authorizationHeader string, parentID int) *NextRoute {
	return &NextRoute{
		Body:                body,
		Query:               query,
		Params:              params,
		AuthorizationHeader: authorizationHeader,
		ParentID:            parentID,
	}
}

func (nrd *NextRoute) GetParentID() int {
	return nrd.ParentID
}

type RouteResponse struct {
	StatusCode int    `json:"statusCode" example:"200"`
	Body       string `json:"body" example:"id=1"`
	ParentID   int    `json:"parentID"`
}

func NewRouteResponse(statusCode, parentID int, body string) *RouteResponse {
	return &RouteResponse{
		StatusCode: statusCode,
		Body:       body,
		ParentID:   parentID,
	}
}

func (rr *RouteResponse) GetParentID() int {
	return rr.ParentID
}
