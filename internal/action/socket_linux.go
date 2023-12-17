package action

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"syscall"
)

var listenConfig = net.ListenConfig{
	Control: func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			if err != nil {
				fmt.Println("Error creating socket: ", err)
			}
		})
	},
}
