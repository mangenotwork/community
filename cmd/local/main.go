package main

import (
	"community/internal/local/registersvc"
	"os"
	"strconv"
)

func main() {

	tag := os.Args[1]

	go registersvc.BurrowClient(str2int(tag))

	//go func() {
	//	//for {
	//	//	name := ""
	//	//	fmt.Scan(&name)
	//	//	b, _ := datasvc.Add(name, name, []byte(name))
	//	//	registersvc.NodeTable.BroadcastData(b)
	//	//}
	//	websvc.Server()
	//}()

	select {}
}

func str2int(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}
