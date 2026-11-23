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

type HTTPError string

const (
	URLError        HTTPError = "URL Error"
	MethodError     HTTPError = "HTTP Method Error"
	ConfigError     HTTPError = "HTTP HandlerConfig Error"
	NewRequestError HTTPError = "NewRequest Error"
	HTTPClientError HTTPError = "HTTP Client Error"
	StatusCodeError HTTPError = "HTTP StatusCode Error"
	BodyError       HTTPError = "HTTP Response Body Error"

	HTTPOK HTTPError = "OK"
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

func (e HTTPError) Error() string {
	return string(e)
}

type HandlerConfig struct {
	URL    string          `json:"url"`
	Method HTTPMethod      `json:"method"`
	Data   json.RawMessage `json:"data"`
}
