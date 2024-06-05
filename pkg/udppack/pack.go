package udppack

import (
	"community/pkg/utils"
	"log"
)

const (
	UDPCodeHeartbeat         UDPCommandCode = 0x00 // 来自节点的心跳报文
	UDPCodeNodeTable         UDPCommandCode = 0x01 // 来自burrow 下发的节点表
	UDPCodeGetNodeTable      UDPCommandCode = 0x02 // 来自节点 节点表请求
	UDPCodeNodeAddr          UDPCommandCode = 0x03 // 来自burrow 返回节点自己的地址
	UDPCodeNodeData          UDPCommandCode = 0x04 // 来自节点 数据
	UDPCodeNodeBroadcastData UDPCommandCode = 0x05 // 来自节点 广播的资料信息
)

func Heartbeat() []byte {
	b, err := PacketEncoder(UDPCodeHeartbeat, []byte("1"))
	if err != nil {
		log.Println(err)
	}
	return b
}

func GetNodeTable(pg int) []byte {
	b, err := PacketEncoder(UDPCodeGetNodeTable, []byte(utils.AnyToString(pg)))
	if err != nil {
		log.Println(err)
	}
	return b
}

func NodeTable(data string) []byte {
	b, err := PacketEncoder(UDPCodeNodeTable, []byte(data))
	if err != nil {
		log.Println(err)
	}
	return b
}

func NodeAddr(data string) []byte {
	b, err := PacketEncoder(UDPCodeNodeAddr, []byte(data))
	if err != nil {
		log.Println(err)
	}
	return b
}

func NodeData(data string) []byte {
	b, err := PacketEncoder(UDPCodeNodeData, []byte(data))
	if err != nil {
		log.Println(err)
	}
	return b
}

func NodeBroadcastData(data []byte) []byte {
	b, err := PacketEncoder(UDPCodeNodeBroadcastData, data)
	if err != nil {
		log.Println(err)
	}
	return b
}
