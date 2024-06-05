package udppack

import (
	"bytes"
	"community/pkg/utils"
	"encoding/binary"
	"fmt"
)

type UDPCommandCode uint8

type Packet struct {
	Code UDPCommandCode // 1个字节
	None []byte         // 7个字节 预留
	Data []byte         // 剩余 540字节 数据内容被压缩
}

func PacketEncoder(code UDPCommandCode, data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, code)
	_ = binary.Write(buf, binary.LittleEndian, []byte("0000000")) // 7个字节 预留
	dCompress := utils.ZlibCompress(data)
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
	bDecompress, err := utils.ZlibDecompress(b)
	if err != nil {
		return nil, err
	}
	return &Packet{
		Code: code,
		Data: bDecompress,
	}, nil
}
