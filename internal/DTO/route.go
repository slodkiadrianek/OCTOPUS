package DTO

import (
	"encoding/json"
	"fmt"
)

type RoutesParentId interface {
	GetParentId() int
}
type CreateRouteData struct {
	Name   string `json:"name"`
	Routes []CreateRoute
}
type CreateRoute struct {
	Path                    string            `json:"path" example:"/users"`
	Method                  string            `json:"method" example:"GET"`
	RequestAuthorization    string            `json:"requestAuthorization" example:"Bearer:fb43fg3487f34g78f3gu"`
	RequestQuery            map[string]string `json:"requestQuery" example:"id=1"`
	RequestParams           map[string]string `json:"requestParams" example:"id=1"`
	RequestBody             map[string]any    `json:"requestBody" example:"id=1"`
	NextRouteBody           []string          `json:"nextRouteBody"`
	NextRouteQuery          []string          `json:"nextRouteQuery"`
	NextRouteParams         []string          `json:"nextRouteParams"`
	NextAuthorizationHeader string            `json:"next_authorization_header"`
	ResponseStatusCode      int               `json:"responseStatusCode" example:"200"`
	ResponseBody            map[string]any    `json:"responseBody" example:"id=1"`
	ParentId                int
}

type JsonMapStringString map[string]string

func (jm *JsonMapStringString) Scan(value interface{}) error {
	if value == nil {
		*jm = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion failed: %T", value)
	}
	return json.Unmarshal(b, jm)
}

type JsonMapStringAny map[string]any

func (ja *JsonMapStringAny) Scan(value interface{}) error {
	if value == nil {
		*ja = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion failed: %T", value)
	}
	return json.Unmarshal(b, ja)
}

type JsonStringSlice []string

func (js *JsonStringSlice) Scan(value interface{}) error {
	if value == nil {
		*js = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion failed: %T", value)
	}
	return json.Unmarshal(b, js)
}

type RouteToTest struct {
	Id                      int
	IpAddress               string
	Port                    string
	Name                    string              `json:"name"`
	Path                    string              `json:"path" example:"/users"`
	Method                  string              `json:"method" example:"GET"`
	RequestAuthorization    string              `json:"requestAuthorization" example:"Bearer:fb43fg3487f34g78f3gu"`
	RequestQuery            JsonMapStringString `json:"requestQuery" example:"id=1"`
	RequestParams           JsonMapStringString `json:"requestParams" example:"id=1"`
	RequestBody             JsonMapStringAny    `json:"requestBody" example:"id=1"`
	NextRouteBody           JsonStringSlice     `json:"nextRouteBody"`
	NextRouteQuery          JsonStringSlice     `json:"nextRouteQuery"`
	NextRouteParams         JsonStringSlice     `json:"nextRouteParams"`
	NextAuthorizationHeader string              `json:"next_authorization_header"`
	ResponseStatusCode      int                 `json:"responseStatusCode" example:"200"`
	ResponseBody            JsonMapStringAny    `json:"responseBody" example:"id=1"`
	ParentId                int
	Status                  string
	AppId                   string
}

type WorkingRoute struct {
	Name            string
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

type NextRouteData struct {
	NextRouteBody           string `json:"nextRouteBody"`
	NextRouteQuery          string `json:"nextRouteQuery"`
	NextRouteParams         string `json:"nextRouteParams"`
	NextAuthorizationHeader string `json:"next_authorization_header"`
	ParentId                int    `json:"parentId"`
}

func NewNextRouteData(nextRouteBody, nextRouteQuery, nextRouteParams, nextAuthorizationHeader string) *NextRouteData {
	return &NextRouteData{
		NextRouteBody:           nextRouteBody,
		NextRouteQuery:          nextRouteQuery,
		NextRouteParams:         nextRouteParams,
		NextAuthorizationHeader: nextAuthorizationHeader,
	}
}

func (nrd *NextRouteData) GetParentId() int {
	return nrd.ParentId
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
