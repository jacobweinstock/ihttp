package file

import (
	"context"
	"io"
	"strings"
	"testing"
)

var data = `#!ipxe
# Boot a persistent RancherOS to RAM

# Location of Kernel/Initrd images
set base-url https://github.com/rancher/k3os/releases/download/v0.7.1/

#chain --autofree https://boot.netboot.xyz

kernel ${base-url}/k3os-vmlinuz-amd64 k3os.password=rancher
initrd ${base-url}/k3os-initrd-amd64
boot
`

func TestScript(t *testing.T) {
	//c := &Config{URI: "/Users/jacobweinstock/ipxe/k3os/k3os.ipxe"}
	c := &Config{RC: io.NopCloser(strings.NewReader(data))}
	o, err := c.Open(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close(context.TODO(), o)
	script, err := c.Read(context.Background(), o)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(string(script))
}
