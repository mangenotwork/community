package main

import (
	"burrow/node"
	"fmt"
	"log"
	"net"
)

var (
	port      = 9981
	nodeTable = node.NewTable()
)

func main() {
	nodeTable.Clear()
	udpServer(port)
}

func udpServer(port int) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		fmt.Println(err)
		return
	}
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
