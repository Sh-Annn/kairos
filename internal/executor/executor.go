package executor

import (
	"context"
	"encoding/json"
	"net/http"
)

type Executor interface {
	Exec(ctx context.Context, config json.RawMessage) ([]byte, ExecError)
}

type ExecError interface {
	Valid() bool
}

type HTTPExecutor struct {
	//	HandlerConfig util.HandlerConfig
	client *http.Client
}

func NewHTTPExecutor() *HTTPExecutor {
	return &HTTPExecutor{
		client: &http.Client{},
	}
}
