//go:build windows

package network

import (
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

func dialIPC(timeout time.Duration) (net.Conn, error) {
	address := `\\.\pipe\wallchemy`
	return winio.DialPipe(address, &timeout)
}
