package registersvc

import (
	"context"
	"fmt"
	"log"
	"net"
)

func NodeTCP(l *net.ListenConfig, port int) {

	listener, err := l.Listen(context.Background(), "tcp", fmt.Sprintf("0.0.0.0:%d", port))

	//listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Println("Error listening:", err)
		return
	}

	log.Println("启动tcp服务...")

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting:", err)
			continue
		}
		go func() {
			defer conn.Close()
			buffer := make([]byte, 512)
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					log.Println("Error reading:", err)
					break
				}
				log.Println("Received:", string(buffer[:n]))
				_, err = conn.Write(buffer[:n])
				if err != nil {
					log.Println("Error writing:", err)
					break
				}
			}
		}()
	}
}
