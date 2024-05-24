package data_block

import (
	"fmt"
	"testing"
)

func Test_Block(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("111")
	bc.AddBlock("222")
	for _, v := range bc.blocks {
		fmt.Printf("Prev hash: %x\n", v.PrevBlockHash)
		fmt.Printf("Data : %s\n", v.Data)
		fmt.Printf("Hash : %x\n", v.Hash)
		fmt.Printf("\n")
	}
}
