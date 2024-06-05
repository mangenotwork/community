package utils

import (
	"bytes"
	"compress/zlib"
	"io"
	"log"
)

func ZlibCompress(src []byte) []byte {
	buf := new(bytes.Buffer)
	writer := zlib.NewWriter(buf)
	_, err := writer.Write(src)
	err = writer.Close()
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

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
