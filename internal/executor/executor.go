package executor

import (
	"context"
	"net/http"
	"time"

	"github.com/Sh-Annn/kairos/util"
)

type Executor interface {
	Exec(ctx context.Context) error
}

type HttpExecutor struct {
	HandlerConfig util.HandlerConfig
	client        *http.Client
}

func NewHttpExecutor(handlerConfig util.HandlerConfig, timeout time.Duration) *HttpExecutor {
	return &HttpExecutor{
		HandlerConfig: handlerConfig,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}
