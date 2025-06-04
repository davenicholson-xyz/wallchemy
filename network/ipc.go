package network

import (
	"fmt"
	"net"
	"runtime"
	"time"
)

func SendIPCMessage(message string) (string, error) {
	var network, address string

	if runtime.GOOS == "windows" {
		network = "pipe"
		address = `\\.\pipe\wallchemy`
	} else {
		network = "unix"
		address = "/tmp/wallchemy.sock"
	}

	conn, err := net.DialTimeout(network, address, 2*time.Second)
	if err != nil {
		return "", fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	if err := conn.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return "", fmt.Errorf("set write deadline failed: %w", err)
	}

	_, err = conn.Write([]byte(message + "\n"))
	if err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return "", fmt.Errorf("set read deadline failed: %w", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("read failed: %w", err)
	}

	return string(buf[:n]), nil
}
