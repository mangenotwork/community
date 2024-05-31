package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "10.0.40.29:9981")
	if err != nil {
		fmt.Println("tcp 连接失败")
		return
	}

	fmt.Printf("tcp 连接成功")

	time.Sleep(3 * time.Second)
	_, err = conn.Write([]byte("来自tcp消息"))
	if err != nil {
		fmt.Println("tcp发送失败:", err)
	}

	time.Sleep(1 * time.Second)
	resp := make([]byte, 256)
	n, err := conn.Read(resp)
	if err != nil {
		fmt.Println("收到tcp恢复失败: ", err)
	}
	fmt.Println("收到tcp回复", string(resp[:n]))
	conn.Close()
}
