package proxy

import (
	"fmt"
	"io"
	"net"
	"time"
)

type Endpoint struct {
	LocalPort  int
	RemoteHost string
	RemotePort int
}

func (e Endpoint) String() string {
	return fmt.Sprintf("%d:%s:%d", e.LocalPort, e.RemoteHost, e.RemotePort)
}

func (e Endpoint) Valid() error {
	if e.LocalPort <= 0 {
		return fmt.Errorf("invalid local port: %d", e.LocalPort)
	}

	if e.RemoteHost == "" {
		return fmt.Errorf("invalid remote host: %s", e.RemoteHost)
	}

	if e.RemotePort <= 0 {
		return fmt.Errorf("invalid remote port: %d", e.RemotePort)
	}

	return nil
}

func (e Endpoint) LocalAddr() string {
	return fmt.Sprintf(":%d", e.LocalPort)
}

func (e Endpoint) RemoteAddr() string {
	return fmt.Sprintf("%s:%d", e.RemoteHost, e.RemotePort)
}

func Handle(conn net.Conn, ep Endpoint) {
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("close conn error: %s\n", err)
		}
	}()

	dialer := net.Dialer{
		Timeout: 3 * time.Second,
	}
	remote, err := dialer.Dial("tcp", ep.RemoteAddr())
	if err != nil {
		fmt.Printf("dial remote: %s\n", err)
		return
	}
	defer func() {
		if cErr := remote.Close(); cErr != nil {
			fmt.Printf("close remote error: %s\n", cErr)
		}
	}()

	go func() {
		n, cErr := io.Copy(remote, conn)
		if cErr != nil {
			fmt.Printf("copy to remote error: %s\n", cErr)
			return
		}
		fmt.Printf("copy to remote %d bytes\n", n)
	}()

	var n int64
	n, err = io.Copy(conn, remote)
	if err != nil {
		fmt.Printf("copy to local error: %s\n", err)
		return
	}

	fmt.Printf("copy from remote %d bytes\n", n)
}
