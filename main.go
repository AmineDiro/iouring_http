package main

import (
	"fmt"
	"github/aminediro/iouring_server/server"
)

func main() {
	l, _ := server.MKRingListener(":8000")
	fmt.Printf("new ring %#v\n", l)

	select {}
}
