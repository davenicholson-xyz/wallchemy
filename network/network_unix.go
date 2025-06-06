//go:build !windows

package network

import (
	"net"
	"time"
)

func dialIPC(timeout time.Duration) (net.Conn, error) {
	address := "/tmp/wallchemy.sock"
	return net.DialTimeout("unix", address, timeout)
}
