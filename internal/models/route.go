package models

import (
	"encoding/json"
	"fmt"
)

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
	ID                      int
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
	ParentID                int
	Status                  string
	AppId                   string
}
