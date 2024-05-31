package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"syscall"
)

func main() {
	go win()
	go win()
	select {}
}

func win() {
	l := &net.ListenConfig{Control: reusePortControl}
	s, err := l.Listen(context.Background(), "tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("启动s ", s.Addr().String())

	defer s.Close()

	for {
		c, err := s.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			// 业务逻辑
		}(c)
	}
}

func reusePortControl(network, address string, c syscall.RawConn) error {

	err := c.Control(func(fd uintptr) {
		// 使用io复用，必须要设置监听socket为SO_REUSEADDR和SO_REUSEPORT

		// 这里是设置地址复用
		err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if err != nil {
			fmt.Println(err)
		}

		////这里是设置端口复用
		//err = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
		//if err != nil {
		//	fmt.Println(err)
		//}

		//这里是设置接收缓冲区大小为10M
		err = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, 10*1024*1024)
		if err != nil {
			fmt.Println(err)
		}

	})
	if err != nil {
		return err
	}
	return nil
}

//func reusePortControlLinux(network, address string, c syscall.RawConn) error {
//
//	err := c.Control(func(fd uintptr) {
//		// syscall.SO_REUSEPORT ,在Linux下还可以指定端口重用
//		err := unix.SetsockoptInt(syscall.Handle(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		err = unix.SetsockoptInt(syscall.Handle(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		err = unix.SetsockoptInt(unix.Handle(fd), unix.SOL_SOCKET, unix.SO_RCVBUF, 10*1024*1024)
//		if err != nil {
//			fmt.Println(err)
//		}
//	})
//	if err != nil {
//		return err
//	}
//	return nil
//}
