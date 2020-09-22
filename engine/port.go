package engine

import (
	"net"
	"strconv"
)

func getPort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	if err := l.Close(); err != nil {
		return "", err
	}

	port := l.Addr().(*net.TCPAddr).Port
	return strconv.Itoa(port), nil
}
