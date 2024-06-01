package node

import (
	"community/pkg/entity"
	"log"
	"net"
	"strings"
	"time"
)

type Table struct {
	List []*entity.NodeTableItem
}

func NewTable() *Table {
	return &Table{
		List: make([]*entity.NodeTableItem, 0),
	}
}

// 获取节点
func (table *Table) Get(addr *net.UDPAddr) string {
	ipList := make([]string, 0)
	for _, v := range table.List {
		if v.Addr.String() != addr.String() {
			ipList = append(ipList, v.Addr.String())
		}
	}
	return strings.Join(ipList, ",")
}

// 设置节点
func (table *Table) Set(addr *net.UDPAddr) {
	log.Println("加入节点", addr.String())
	has := false
	for _, v := range table.List {
		if v.Addr.String() == addr.String() {
			has = true
			v.NextTime = time.Now().Unix()
		}
	}
	if !has {
		table.List = append(table.List, entity.NewNodeTableItem(addr))
	}
}

// 删除节点
func (table *Table) del(addr *net.UDPAddr) {
	for i, v := range table.List {
		if v.Addr.String() == addr.String() {
			table.List = append(table.List[:i], table.List[i+1:]...)
		}
	}
}

// 定期清理心跳过期的节点
func (table *Table) Clear() {
	go func() {
		ticker := time.NewTicker(time.Second * 10) // 10秒执行一次
		for {
			select {
			case <-ticker.C:
				// 这里执行你的定时任务
				log.Println("定时任务运行中...")
				now := time.Now().Unix()
				for _, v := range table.List {
					log.Println("节点:", v.Addr, "   time:", v.NextTime)
					if now-v.NextTime > 10 {
						log.Println("清理心跳过期的节点 -> ", v.Addr.String())
						table.del(v.Addr)
					}
				}
			}
		}
	}()
}
