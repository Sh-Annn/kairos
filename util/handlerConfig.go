package util

import (
	"encoding/json"
)

type HTTPMethod string

const (
	MethodGet    HTTPMethod = "GET"
	MethodPost   HTTPMethod = "POST"
	MethodPut    HTTPMethod = "PUT"
	MethodDelete HTTPMethod = "DELETE"
)

func (m HTTPMethod) Valid() bool {
	switch m {
	case MethodGet, MethodPost, MethodPut, MethodDelete:
		return true
	default:
		return false
	}
}

type HandlerConfig struct {
	URL    string          `json:"url"`
	Method HTTPMethod      `json:"method"`
	Data   json.RawMessage `json:"data"`
}
