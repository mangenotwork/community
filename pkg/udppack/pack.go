package udppack

import (
	"community/pkg/utils"
	"log"
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
