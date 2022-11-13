package main

import (
	"fmt"

	"github.com/aminediro/iouring_server/server"
)

func main() {
	l, _ := server.MKRingListener(":8000")
	fmt.Printf("new ring %#v\n", l)

	select {}
}
