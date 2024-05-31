package main

import (
	"local/registersvc"
	"os"
	"strconv"
)

func main() {

	tag := os.Args[1]

	registersvc.BurrowClient(str2int(tag))

	select {}
}

func str2int(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}
