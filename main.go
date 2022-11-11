package main

//#cgo LDFLAGS: -luring
//#include "server.h"
//#include <liburing.h>
import "C"

import "fmt"

func main() {
	r := RingInit()
	fmt.Println("new ring", r)
}

func RingInit() int {
	r := int(C.queue_init())
	return r
}
