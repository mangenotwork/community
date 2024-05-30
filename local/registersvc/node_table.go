package registersvc

import (
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
		port, _ := strconv.Atoi(t[1])
		newTable = append(newTable, &Node{
			Addr: v,
			IP:   t[0],
			Port: port,
		})
	}
	table.List = newTable
}

func (table *Table) Hello() {
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

					conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
					if err != nil {
						fmt.Println(err)
						continue
					}

					// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
					if _, err = conn.Write([]byte(fmt.Sprintf("hello, my is %s", srcAddr.String()))); err != nil {
						log.Println("send handshake:", err)
					}

					_ = conn.Close()
				}

			}
		}
	}()
}
