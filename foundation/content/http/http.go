package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	RC io.ReadCloser
}

func (c *Config) Open(ctx context.Context, uri string) (io.ReadCloser, error) {
	if c.RC == nil {
		req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
		if err != nil {
			return nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode < 200 && resp.StatusCode >= 300 {
			resp.Body.Close()
			return nil, fmt.Errorf(resp.Status)
		}
		return resp.Body, nil
	}
	return c.RC, nil
}

func (c *Config) Close(_ context.Context, rc io.ReadCloser) error {
	return rc.Close()
}

func (c *Config) Read(ctx context.Context, r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
