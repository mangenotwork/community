package main

import (
	"community/internal/burrow/node"
	"community/pkg/conf"
	"context"
	"fmt"
	"log"
	"net"
	"syscall"
)

var (
	nodeTable = node.NewTable()
)

func main() {
	conf.InitBurrowConf()
	nodeTable.Clear()
	udpServer(conf.BurrowConf.Port)
}

func udpServer(port int) {
	l := &net.ListenConfig{Control: reusePortControl}
	lp, err := l.ListenPacket(context.Background(), "udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Println(err)
	}

	go func() {
		// test tcp 端口复用
		ltcp, err := l.Listen(context.Background(), "tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			log.Printf("Could not start TCP listener: %s", err)
			return
		}
		log.Println("tcp服务启动...")
		for {
			c, err := ltcp.Accept()
			if err != nil {
				log.Printf("Listener returned: %s", err)
				break
			}
			go handleConnection(c)
		}
	}()

	go func() {
		// test tcp 端口复用
		ltcp, err := l.Listen(context.Background(), "tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			log.Printf("Could not start TCP listener: %s", err)
			return
		}
		log.Println("tcp服务启动...")
		for {
			c, err := ltcp.Accept()
			if err != nil {
				log.Printf("Listener returned: %s", err)
				break
			}
			go handleConnection(c)
		}
	}()

	//listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	listener := lp.(*net.UDPConn)

	log.Printf("本地地址: <%s> \n", listener.LocalAddr().String())

	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		nodeTable.Set(remoteAddr)

		_, err = listener.WriteToUDP([]byte("2"+remoteAddr.String()), remoteAddr)
		if err != nil {
			log.Println("下发节点地址失败: ", err)
		}

		// 下发节点表
		nodeData := nodeTable.Get(remoteAddr)
		if len(nodeData) > 0 {
			_, wErr := listener.WriteToUDP([]byte("0"+nodeData), remoteAddr)
			if wErr != nil {
				log.Println("下发节点表Err: ", wErr)
			}
		}

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

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		fmt.Println("Received:", string(buffer[:n]))

		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}
}
