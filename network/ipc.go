package network

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/Microsoft/go-winio" // Add this import
)

func SendIPCMessage(message string) (string, error) {
	var conn net.Conn
	var err error
	timeout := 2 * time.Second

	if runtime.GOOS == "windows" {
		// Use winio for Windows named pipes
		address := `\\.\pipe\wallchemy`
		conn, err = winio.DialPipe(address, &timeout)
	} else {
		// Unix domain socket
		address := "/tmp/wallchemy.sock"
		conn, err = net.DialTimeout("unix", address, timeout)
	}

	if err != nil {
		return "", fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return "", fmt.Errorf("set write deadline failed: %w", err)
	}

	// Send message with newline terminator
	if _, err := conn.Write([]byte(message + "\n")); err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}

	// Set read deadline
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return "", fmt.Errorf("set read deadline failed: %w", err)
	}

	// Read response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("read failed: %w", err)
	}

	return string(buf[:n]), nil
}
