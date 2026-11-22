package executor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Sh-Annn/kairos/internal/domain"
)

func (h *HTTPExecutor) Exec(ctx context.Context, config json.RawMessage) ([]byte, ExecError) {
	var handlerConfig domain.HandlerConfig
	err := json.Unmarshal(config, &handlerConfig)
	if err != nil {
		return nil, domain.ConfigError
	}
	// TODO: diffrent error

	if !handlerConfig.Method.Valid() {
		return nil, domain.MethodError
	}

	if handlerConfig.URL == "" {
		return nil, domain.URLError
	}

	req, err := http.NewRequestWithContext(ctx, string(handlerConfig.Method), handlerConfig.URL, nil)
	if err != nil {
		return nil, domain.NewRequestError
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, domain.HTTPClientError
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, domain.StatusCodeError
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return nil, domain.BodyError
	}
	return body, domain.HTTPOK
}
