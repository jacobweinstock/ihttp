package file

import (
	"context"
	"io"
	"os"
)

type Config struct {
	RC  io.ReadCloser
}

func (c *Config) Open(_ context.Context, uri string) (io.ReadCloser, error) {
	if c.RC == nil {
		return os.Open(uri)
	}

	return c.RC, nil
}

func (c *Config) Close(_ context.Context, rc io.ReadCloser) error {
	return rc.Close()
}

func (c *Config) Read(ctx context.Context, r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
