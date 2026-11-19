package executor

import (
	"context"
	"net/http"

	"github.com/Sh-Annn/kairos/internal/domain"
)

type Executor interface {
	Exec(ctx context.Context, handlerConfig domain.HandlerConfig) ([]byte, error)
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
