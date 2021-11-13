package file

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	URI string
	RC  io.ReadCloser
}

type InMemoryYamlFile struct {
	Machines []Machine `yaml:"machines"`
}

type Machine struct {
	Mac    string `yaml:"mac"`
	IP     string `yaml:"ip"`
	Script string `yaml:"script"`
}

func (c *Config) Open(ctx context.Context) (io.ReadCloser, error) {
	if c.RC == nil {
		return os.Open(c.URI)
	}
	return c.RC, nil
}

func (c *Config) Close(_ context.Context, r io.ReadCloser) error {
	return r.Close()
}

func (c *Config) Locate(ctx context.Context, mac net.HardwareAddr, r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	var b InMemoryYamlFile

	if err = yaml.Unmarshal(data, &b); err != nil {
		return "", err
	}

	for _, m := range b.Machines {
		if m.Mac == mac.String() {
			return m.Script, nil
		}
	}

	return "", fmt.Errorf("no file found for mac %s", mac.String())
}
