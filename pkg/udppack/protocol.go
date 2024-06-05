package udppack

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type UDPCommandCode uint8

const (
	UDPCodeHeartbeat    UDPCommandCode = 0x00 // 来自节点的心跳报文
	UDPCodeNodeTable    UDPCommandCode = 0x01 // 来自burrow 下发的节点表
	UDPCodeGetNodeTable UDPCommandCode = 0x02 // 来自节点 节点表请求
	UDPCodeNodeAddr     UDPCommandCode = 0x03 // 来自burrow 返回节点自己的地址
	UDPCodeNodeData     UDPCommandCode = 0x04 // 来自节点 数据
)

type Packet struct {
	Code UDPCommandCode // 1个字节
	None []byte         // 7个字节 预留
	Data []byte         // 剩余 540字节 数据内容被压缩
}

func PacketEncoder(code UDPCommandCode, data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, code)
	_ = binary.Write(buf, binary.LittleEndian, []byte("0000000")) // 7个字节 预留
	dCompress := ZlibCompress(data)
	_ = binary.Write(buf, binary.LittleEndian, dCompress)
	b := buf.Bytes()
	// 不允许大包传输
	if len(b) > 548 {
		return nil, fmt.Errorf("数据太大请用TCP包")
	}
	return b, nil
}

func PacketDecrypt(data []byte, n int) (*Packet, error) {
	var err error
	if n < 8 {
		return nil, fmt.Errorf("空包")
	}
	code := UDPCommandCode(data[0:1][0])
	b := data[8:n]
	bDecompress, err := ZlibDecompress(b)
	if err != nil {
		return nil, err
	}
	return &Packet{
		Code: code,
		Data: bDecompress,
	}, nil
}

// ZlibCompress zlib压缩
func ZlibCompress(src []byte) []byte {
	buf := new(bytes.Buffer)
	//根据创建的buffer生成 zlib writer
	writer := zlib.NewWriter(buf)
	//写入数据
	_, err := writer.Write(src)
	err = writer.Close()
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

// ZlibDecompress zlib解压
func ZlibDecompress(src []byte) ([]byte, error) {
	reader := bytes.NewReader(src)
	gr, err := zlib.NewReader(reader)
	if err != nil {
		return []byte(""), err
	}
	bf := make([]byte, 0)
	buf := bytes.NewBuffer(bf)
	_, err = io.Copy(buf, gr)
	err = gr.Close()
	return buf.Bytes(), err
}
