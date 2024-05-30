package registersvc

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var srcAddr *net.UDPAddr

func BurrowClient(port int) {
	// 当前进程标记字符串,便于显示
	srcAddr = &net.UDPAddr{IP: net.IPv4zero, Port: port} // 注意端口必须固定
	dstAddr := &net.UDPAddr{IP: net.ParseIP("10.0.40.29"), Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	go func() {
		for {
			data := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(data)
			if err != nil {
				fmt.Printf("error during read: %s", err)
			}
			log.Println("收到节点表:", string(data[:n]))
			NodeTable.Refresh(string(data[:n]))
		}
	}()

	for {
		if _, err = conn.Write([]byte("心跳")); err != nil {
			log.Panic(err)
		}
		time.Sleep(5 * time.Second) // 5秒发一次心跳包
	}

	//anotherPeer := parseAddr(string(data[:n]))
	//fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, anotherPeer.String())

	//// 开始打洞
	//bidirectionHole(srcAddr, &anotherPeer)
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, err = conn.Write([]byte("hello")); err != nil {
		log.Println("send handshake:", err)
	}

	go func() {
		time.Sleep(13 * time.Second)
		tcpc(anotherAddr)
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from []")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()

	go func() {
		tcps()
	}()

	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("收到数据:%s\n", data[:n])
		}
	}
}

func tcpc(addr *net.UDPAddr) {
	ser := addr.String()
	log.Println(ser)

	conn, err := net.Dial("tcp", ser)
	if err != nil {
		fmt.Println("tcp 连接失败")
		return
	}
	defer conn.Close()

	go func() {
		for {
			resp := make([]byte, 256)
			n, err := conn.Read(resp)
			if err != nil {
				fmt.Println("收到tcp恢复失败: ", err)
			}
			fmt.Println("收到tcp回复", string(resp[:n]))
		}
	}()

	fmt.Printf("tcp 连接成功")
	for {
		time.Sleep(3 * time.Second)
		_, err = conn.Write([]byte("来自tcp消息"))
		if err != nil {
			fmt.Println("tcp发送失败:", err)
		}
	}
}

func tcps() {
	listener, err := net.Listen("tcp", ":9982")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("启动tcp服务...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go func() {
			defer conn.Close()
			buffer := make([]byte, 512)
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Error reading:", err)
					break
				}
				fmt.Println("Received:", string(buffer[:n]))
				_, err = conn.Write(buffer[:n])
				if err != nil {
					fmt.Println("Error writing:", err)
					break
				}
			}
		}()
	}
}
