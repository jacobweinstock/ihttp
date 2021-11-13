package file

import (
	"context"
	"net"
	"strings"
	"testing"
)

var data = `---
machines:
- mac: 08:00:27:29:4e:67
  ip: 192.168.2.3
  script: ~/ipxe/k3os/k3os.ipxe
- mac: 08:00:27:9e:ec:94
  ip: 192.168.2.4
- mac: 08:00:27:b1:da:42
  ip: 192.168.2.5
- mac: d2:9a:c5:cb:46:ea
  ip: 192.168.2.6
  script: ~/ipxe/k3os/k3os.ipxe
`

func TestLocation(t *testing.T) {
	r := strings.NewReader(data)
	f := Config{}
	location, err := f.Locate(context.Background(), net.HardwareAddr{0x08, 0x00, 0x27, 0x29, 0x4e, 0x67}, r)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(location)
}
