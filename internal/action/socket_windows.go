package action

import (
	"fmt"
	"golang.org/x/sys/windows"
	"net"
	"syscall"
)

var listenConfig = net.ListenConfig{
	Control: func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			err := windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
			if err != nil {
				fmt.Println("Error creating socket: ", err)
			}
		})
	},
}
