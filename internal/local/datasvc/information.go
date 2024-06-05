package datasvc

import (
	"bufio"
	"community/pkg/db"
	"community/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// 资料主体
type Information struct {
	Id       string      // 唯一id
	Title    string      // 标题
	Describe string      // 描述-简介
	Time     int64       // 发布时间
	Content  []*DataItem // 资料内容 - 分片数据文件
}

// 分片数据
type DataItem struct {
	Hash string // 唯一标识
	Path string // 文件位置
}

// 内存缓存存储表
// 异步存储到磁盘
// 写:  内存缓存  -> 磁盘
// 读:  内存缓存（没有 <- 磁盘）
// 删:  游标位置，长度

var InformationTable = sync.Map{}

// 创建数据
func Add(title, describe string, data []byte) ([]byte, error) {
	id := utils.IDMd5()
	idPath := fmt.Sprintf("./data/%s", id)
	_ = os.MkdirAll(idPath, 0666)

	// todo 对数据进行分片存储到文件
	// 需要进行压缩, 分片单文件再1M左右
	dataHash := utils.IDMd5()
	dataPath := fmt.Sprintf("./data/%s/%s", id, dataHash)

	fileHandle, err := os.OpenFile(dataPath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("open file error :", err)
		return []byte(""), err
	}
	defer fileHandle.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
	buf := bufio.NewWriter(fileHandle)
	_, _ = buf.Write(data)
	err = buf.Flush()
	if err != nil {
		log.Println("flush error :", err)
	}

	info := Information{
		Id:       utils.IDMd5(),
		Title:    title,
		Describe: describe,
		Time:     time.Now().UnixMilli(),
		Content: []*DataItem{
			{
				Hash: dataHash,
			},
		},
	}
	InformationTable.Store(info.Id, info)
	b, _ := json.Marshal(info)
	_ = db.InformationTable.Set(info.Id, b)

	return b, nil
}

func AddFromId(data *Information) {
	InformationTable.Store(data.Id, data)
	b, _ := json.Marshal(data)
	_ = db.InformationTable.Set(data.Id, b)
}
