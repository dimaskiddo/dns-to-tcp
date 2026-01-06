//go:build linux || darwin

package main

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func setSocketOptions(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_RCVBUF, 4*1024*1024)
		unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
		unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	})
}
