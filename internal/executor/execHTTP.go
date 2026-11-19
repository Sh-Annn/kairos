package executor

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (h *HttpExecutor) Exec(ctx context.Context) error {
	if !h.HandlerConfig.Method.Valid() {
		return fmt.Errorf("HTTP Method is not supported.")
	}

	req, err := http.NewRequestWithContext(ctx, string(h.HandlerConfig.Method), h.HandlerConfig.URL, nil)
	if err != nil {
		log.Println(fmt.Errorf("cannot create request: %w", err))
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	log.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad http statusCode")
	}
	log.Println(resp.Header)

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(body))
	return nil
}
