package registersvc

import (
	"community/pkg/logger"
	"community/pkg/udppack"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var NodeTable = &Table{List: make([]*Node, 0)}

type Table struct {
	List []*Node
}

type Node struct {
	Addr string
	IP   string
	Port int
}

func (table *Table) Refresh(listStr string) {
	newTable := make([]*Node, 0)
	for _, v := range strings.Split(listStr, ",") {
		if v == "" {
			continue
		}
		t := strings.Split(v, ":")
		log.Println("解析节点: ", t)
		port, _ := strconv.Atoi(t[1])
		newTable = append(newTable, &Node{
			Addr: v,
			IP:   t[0],
			Port: port,
		})
	}
	table.List = newTable
}

func (table *Table) Hello(listen *net.UDPConn) {
	go func() {
		ticker := time.NewTicker(time.Second * 10) // 10秒执行一次
		for {
			select {

			case <-ticker.C:
				for _, v := range table.List {
					anotherAddr := &net.UDPAddr{
						IP:   net.ParseIP(v.IP),
						Port: v.Port,
					}
					if _, err := listen.WriteTo(udppack.NodeData("hello, my is "+myAddr), anotherAddr); err != nil {
						logger.Error("send handshake:", err)
					}
				}

			}
		}
	}()
}

func (table *Table) HelloTcp() {
	go func() {
		ticker := time.NewTicker(time.Second * 10) // 10秒执行一次
		for {
			select {

			case <-ticker.C:
				for _, v := range table.List {
					anotherAddr := &net.UDPAddr{
						IP:   net.ParseIP(v.IP),
						Port: v.Port,
					}
					tcpc(anotherAddr)
				}

			}
		}
	}()
}

func tcpc(addr *net.UDPAddr) {
	ser := addr.String()
	log.Println(ser)

	conn, err := net.Dial("tcp", ser)
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
