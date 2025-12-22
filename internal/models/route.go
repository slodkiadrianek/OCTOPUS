package models

import (
	"encoding/json"
	"fmt"
)

type JSONMapStringString map[string]string

func (jm *JSONMapStringString) Scan(value any) error {
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

type JSONMapStringAny map[string]any

func (ja *JSONMapStringAny) Scan(value any) error {
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

type JSONStringSlice []string

func (js *JSONStringSlice) Scan(value any) error {
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
	IPAddress               string
	Port                    string
	Name                    string              `json:"name"`
	Path                    string              `json:"path" example:"/users"`
	Method                  string              `json:"method" example:"GET"`
	RequestAuthorization    string              `json:"requestAuthorization" example:"Bearer:fb43fg3487f34g78f3gu"`
	RequestQuery            JSONMapStringString `json:"requestQuery" example:"id=1"`
	RequestParams           JSONMapStringString `json:"requestParams" example:"id=1"`
	RequestBody             JSONMapStringAny    `json:"requestBody" example:"id=1"`
	NextRouteBody           JSONStringSlice     `json:"nextRouteBody"`
	NextRouteQuery          JSONStringSlice     `json:"nextRouteQuery"`
	NextRouteParams         JSONStringSlice     `json:"nextRouteParams"`
	NextAuthorizationHeader string              `json:"next_authorization_header"`
	ResponseStatusCode      int                 `json:"responseStatusCode" example:"200"`
	ResponseBody            JSONMapStringAny    `json:"responseBody" example:"id=1"`
	ParentID                int
	Status                  string
	AppID                   string
}
