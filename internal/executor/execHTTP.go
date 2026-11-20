package executor

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Sh-Annn/kairos/internal/domain"
)

func (h *HTTPExecutor) Exec(ctx context.Context, handlerConfig domain.HandlerConfig) ([]byte, error) {
	// TODO: diffrent error

	if !handlerConfig.Method.Valid() {
		return nil, fmt.Errorf(string(domain.MethodError))
	}

	if handlerConfig.URL == "" {
		return nil, fmt.Errorf(string(domain.URLError))
	}

	req, err := http.NewRequestWithContext(ctx, string(handlerConfig.Method), handlerConfig.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client cannot call http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bad http statusCode: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return nil, fmt.Errorf("cannot read from response body: %w", err)
	}
	return body, nil
}
