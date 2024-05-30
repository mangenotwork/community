package node

import (
	"net"
	"time"
)

type Info struct {
	Addr     *net.UDPAddr
	NextTime int64 // 时间戳秒
}

func NewInfo(addr *net.UDPAddr) *Info {
	return &Info{
		Addr:     addr,
		NextTime: time.Now().Unix(),
	}
}
