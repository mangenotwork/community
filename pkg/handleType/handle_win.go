//go:build windows

package handleType

import (
	"fmt"
	"syscall"
)

func ReusePortControl(network, address string, c syscall.RawConn) error {

	err := c.Control(func(fd uintptr) {

		// 使用io复用，必须要设置监听socket为SO_REUSEADDR和SO_REUSEPORT

		// 这里是设置地址复用
		err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if err != nil {
			fmt.Println(err)
		}

		////这里是设置接收缓冲区大小为10M
		//err = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, 10*1024*1024)
		//if err != nil {
		//	fmt.Println(err)
		//}

	})
	if err != nil {
		return err
	}
	return nil
}
