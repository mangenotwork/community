package main

import (
	"fmt"
	"net"
	"strings"
)

// 记录来自节点的心跳，并记录节点信息
var nodeMap []string

// 存储节点信息
func addNode(addr string) {
	has := false
	for _, v := range nodeMap {
		if v == addr {
			has = true
		}
	}
	if !has {
		nodeMap = append(nodeMap, addr)
	}
}

func nodeInfo() string {
	return strings.Join(nodeMap, ",")
}

func main() {
	// 监听
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9981})
	if err != nil {
		fmt.Println(err)
		return
	}
	nodeAddr := listener.LocalAddr().String()
	fmt.Printf("Local: <%s> \n", nodeAddr)

	// 读取数据
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
			continue
		}
		// 将节点加入到节点表
		addNode(remoteAddr.String())

		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		_, err = listener.WriteToUDP([]byte(nodeInfo()), remoteAddr)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}

}
