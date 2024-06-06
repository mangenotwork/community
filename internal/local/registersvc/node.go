package registersvc

import (
	"community/internal/local/datasvc"
	"community/pkg/logger"
	"community/pkg/udppack"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"syscall"
	"time"
)

var listen *net.UDPConn

var myAddr string

func reusePortControl(network, address string, c syscall.RawConn) error {
	return c.Control(func(fd uintptr) {
		err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if err != nil {
			fmt.Println(err)
		}
	})
}

func BurrowClient(port int) {

	var err error

	l := &net.ListenConfig{Control: reusePortControl}

	// 启动tcp
	go NodeTCP(l, port)

	// 启动udp
	lp, err := l.ListenPacket(context.Background(), "udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger.Error(err)
	}

	//srcAddr := &net.UDPAddr{
	//	IP:   net.IPv4zero,
	//	Port: port,
	//}
	//
	//listen, err = net.ListenUDP("udp", srcAddr)

	listen = lp.(*net.UDPConn)

	if err != nil {
		logger.Error("启动UDP失败")
		panic(err)
	}
	logger.Info("启动UDP服务...")

	// burrow地址
	burrow := "110.41.4.212" // "10.0.40.29"
	dstAddr := &net.UDPAddr{IP: net.ParseIP(burrow), Port: 9981}

	NodeTable.Hello(listen)
	NodeTable.HelloTcp()

	// 接收数据
	go func() {
		for {
			data := make([]byte, 1024)
			n, _, err := listen.ReadFromUDP(data)
			if err != nil {
				logger.ErrorF("error during read: %s", err.Error())
			}

			pack, err := udppack.PacketDecrypt(data[:n], n)
			if err != nil {
				logger.Error(err)
				continue
			}

			switch pack.Code {

			case udppack.UDPCodeNodeAddr:
				myAddr = string(pack.Data)
				logger.Info("myAddr = ", myAddr)
				// 发送节点表请求
				_, err = listen.WriteTo(udppack.GetNodeTable(0), dstAddr)

			case udppack.UDPCodeNodeTable:
				tableData := string(pack.Data)
				logger.Info("收到节点表 = ", tableData)
				NodeTable.Refresh(tableData)

			case udppack.UDPCodeNodeBroadcastData:
				information := string(pack.Data)
				logger.Info("收到节点广播的数据 = ", information)
				newInformation := &datasvc.Information{}
				_ = json.Unmarshal(pack.Data, &newInformation)
				datasvc.AddFromId(newInformation)

			default:
				logger.Info("收到来自节点的数据:", string(pack.Data))
			}

			//dataStr := string(data[:n])
			//if len(dataStr) < 1 {
			//	continue
			//}
			//
			//switch string(dataStr[0]) {
			//case "0":
			//	nodeTableData := dataStr[1:len(dataStr)]
			//	log.Println("收到节点表:", nodeTableData)
			//	NodeTable.Refresh(nodeTableData)
			//case "2":
			//	myAddr = dataStr[1:len(dataStr)]
			//default:
			//	msg := dataStr[1:len(dataStr)]
			//	log.Println("收到来自节点的数据:", msg)
			//}

		}
	}()

	for {
		logger.Info("发送心跳包")
		_, err = listen.WriteTo(udppack.Heartbeat(), dstAddr)
		if err != nil {
			logger.Error("发生心跳包失败")
		}
		time.Sleep(5 * time.Second) // 5秒发一次心跳包
	}

}
