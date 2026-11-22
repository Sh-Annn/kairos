package domain

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

// TODO: define HTTP Error

func (m HTTPMethod) Valid() bool {
	switch m {
	case MethodGet, MethodPost, MethodPut, MethodDelete:
		return true
	default:
		return false
	}
}

type HTTPError int

const (
	URLError HTTPError = iota
	MethodError
	ConfigError
	NewRequestError
	HTTPClientError
	StatusCodeError
	BodyError

	HTTPOK
)

func (e HTTPError) Valid() bool {
	switch e {
	case URLError, MethodError, ConfigError, HTTPClientError,
		NewRequestError, StatusCodeError, BodyError:
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
