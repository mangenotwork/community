package entity

import (
	"net"
	"time"
)

type NodeTableItem struct {
	Addr     *net.UDPAddr
	NextTime int64 // 时间戳秒
}

func NewNodeTableItem(addr *net.UDPAddr) *NodeTableItem {
	return &NodeTableItem{
		Addr:     addr,
		NextTime: time.Now().Unix(),
	}
}
