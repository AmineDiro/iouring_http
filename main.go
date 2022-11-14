package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aminediro/iouring_server/server"
)

func main() {
	server.IncreaseResources()

	l, err := server.MKRingListener(":8000")
	if err != nil {
		panic(err)
	}
	fmt.Printf("new ring %#v\n", l)

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	select {}
}
