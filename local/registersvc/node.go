package registersvc

import (
	"fmt"
	"log"
	"net"
	"time"
)

var listen *net.UDPConn
var err error
var myAddr string

func BurrowClient(port int) {

	srcAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	}

	listen, err = net.ListenUDP("udp", srcAddr)
	if err != nil {
		log.Println("启动UDP失败")
		panic(err)
	}
	log.Println("启动UDP服务...")

	// 启动tcp
	go NodeTCP(port)

	// burrow地址
	dstAddr := &net.UDPAddr{IP: net.ParseIP("10.0.40.29"), Port: 9981}

	NodeTable.Hello(listen)
	NodeTable.HelloTcp()

	// 接收数据
	go func() {
		for {
			data := make([]byte, 1024)
			n, _, err := listen.ReadFromUDP(data)
			if err != nil {
				fmt.Printf("error during read: %s", err)
			}

			dataStr := string(data[:n])
			if len(dataStr) < 1 {
				continue
			}
			switch string(dataStr[0]) {
			case "0":
				nodeTableData := dataStr[1:len(dataStr)]
				log.Println("收到节点表:", nodeTableData)
				NodeTable.Refresh(nodeTableData)
			case "2":
				myAddr = dataStr[1:len(dataStr)]
			default:
				msg := dataStr[1:len(dataStr)]
				log.Println("收到来自节点的数据:", msg)
			}

		}
	}()

	for {
		log.Println("发生心跳包保活...")
		_, err = listen.WriteTo([]byte("心跳"), dstAddr)
		if err != nil {
			log.Println("发生心跳包失败")
		}
		time.Sleep(5 * time.Second) // 5秒发一次心跳包
	}

}
