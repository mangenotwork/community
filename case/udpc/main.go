package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var my = "127.0.0.1:12312"
var srcAddr = &net.UDPAddr{IP: net.IPv4zero, Port: 12312}

var nodTable = make([]string, 0)

var linkMap = make(map[string]int)

func main() {

	go register()

	go broadcast()

	select {}
}

// 在 burrow服务注册节点信息，并更新自己的节点表
func register() {
	sip := net.ParseIP("127.0.0.1")

	dstAddr := &net.UDPAddr{IP: sip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	my = conn.LocalAddr().String()
	fmt.Println("我自己的地址: ", my)

	go func() {
		data := make([]byte, 1024)
		for {
			n, remoteAddr, err := conn.ReadFromUDP(data)
			if err != nil {
				fmt.Printf("error during read: %s", err)
			}
			fmt.Printf("<%s> %s\n", remoteAddr, data[:n])

			// 更新节点表
			nodTable = strings.Split(string(data[:n]), ",")
		}
	}()

	for {
		time.Sleep(3 * time.Second)
		conn.Write([]byte("hello"))
		//fmt.Printf("<%s>\n", conn.RemoteAddr())
	}
}

// 广播消息
func broadcast() {
	for {
		time.Sleep(2 * time.Second)
		if len(nodTable) == 0 {
			continue
		}
		for _, v := range nodTable {
			if my == v {
				continue
			}
			if _, has := linkMap[v]; has {
				continue
			}

			go func(v string) {
				linkMap[v] = 1
				addr := strings.Split(v, ":")
				sip := net.ParseIP(addr[0])
				dstAddr := &net.UDPAddr{IP: sip, Port: str2int(addr[1])}
				conn, err := net.DialUDP("udp", srcAddr, dstAddr)
				if err != nil {
					fmt.Println(err)
				}
				defer conn.Close()

				go func() {
					data := make([]byte, 1024)
					for {
						n, _, err := conn.ReadFromUDP(data)
						if err != nil {
							fmt.Printf("error during read: %s", err)
						}
						fmt.Printf("收到节点广播消息 -> %s\n", data[:n])
					}
				}()

				for {
					time.Sleep(5 * time.Second)
					conn.Write([]byte("我是节点 ： " + my))
					//fmt.Printf("<%s>\n", conn.RemoteAddr())
				}
			}(v)

		}
	}

}

func str2int(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}
