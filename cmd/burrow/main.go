package main

import (
	"community/internal/burrow/node"
	"community/pkg/conf"
	"community/pkg/handleType"
	"community/pkg/logger"
	"community/pkg/udppack"
	"context"
	"fmt"
	"net"
)

var (
	nodeTable = node.NewTable()
)

func main() {
	conf.InitBurrowConf()
	nodeTable.Clear()
	udpServer(conf.BurrowConf.Port)
}

func udpServer(port int) {
	l := &net.ListenConfig{Control: handleType.ReusePortControl}
	lp, err := l.ListenPacket(context.Background(), "udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger.Error(err)
	}

	go func() {
		// test tcp 端口复用
		lTcp, err := l.Listen(context.Background(), "tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			logger.Info("Could not start TCP listener: %s", err)
			return
		}
		logger.Info("tcp服务启动...")
		for {
			c, err := lTcp.Accept()
			if err != nil {
				logger.ErrorF("Listener returned: %s", err.Error())
				break
			}
			go handleConnection(c)
		}
	}()

	go func() {
		// test tcp 端口复用
		lTcp, err := l.Listen(context.Background(), "tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			logger.ErrorF("Could not start TCP listener: %s", err.Error())
			return
		}
		logger.Info("tcp服务启动...")
		for {
			c, err := lTcp.Accept()
			if err != nil {
				logger.ErrorF("Listener returned: %s", err.Error())
				break
			}
			go handleConnection(c)
		}
	}()

	//listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	listener := lp.(*net.UDPConn)

	logger.InfoF("本地地址: <%s> \n", listener.LocalAddr().String())

	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			logger.ErrorF("error during read: %s", err.Error())
		}

		pack, err := udppack.PacketDecrypt(data[:n], n)
		if err != nil {
			logger.Error(err)
			continue
		}

		logger.InfoF("<%s> %s\n", remoteAddr.String(), data[:n])
		nodeTable.Set(remoteAddr)

		switch pack.Code {

		case udppack.UDPCodeHeartbeat:
			_, err = listener.WriteToUDP(udppack.NodeAddr(remoteAddr.String()), remoteAddr)
			if err != nil {
				logger.Error("下发节点地址失败: ", err)
			}

		case udppack.UDPCodeGetNodeTable:
			logger.Info("下发节点表")
			nodeData := nodeTable.Get(remoteAddr)
			if len(nodeData) > 0 {
				_, wErr := listener.WriteToUDP(udppack.NodeTable(nodeData), remoteAddr)
				if wErr != nil {
					logger.Error("下发节点表Err: ", wErr)
				}
			}

		}

		//_, err = listener.WriteToUDP([]byte("2"+remoteAddr.String()), remoteAddr)
		//if err != nil {
		//	log.Println("下发节点地址失败: ", err)
		//}

		//// 下发节点表
		//nodeData := nodeTable.Get(remoteAddr)
		//if len(nodeData) > 0 {
		//	_, wErr := listener.WriteToUDP([]byte("0"+nodeData), remoteAddr)
		//	if wErr != nil {
		//		log.Println("下发节点表Err: ", wErr)
		//	}
		//}

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			logger.Error("Error reading:", err)
			return
		}
		logger.Info("Received:", string(buffer[:n]))

		_, err = conn.Write(buffer[:n])
		if err != nil {
			logger.Error("Error writing:", err)
			return
		}
	}
}
