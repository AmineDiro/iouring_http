package main

import (
	"fmt"
	"github/aminediro/iouring_server/server"
)

func main() {
	l, _ := server.MKRingListener()
	fmt.Println("new ring", l)

	// HTTP Server using Liburing

	// Create a Ring

	// Create a Socket(filedesc)
	// Setup a ring-based Socker listener (fd)
	/// Should implement : Accept() - Addr() - Close()
	///

	// Pass the Listener to an HttpMux Server
}
