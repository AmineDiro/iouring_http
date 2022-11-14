package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"os"
	"os/signal"
	"syscall"

	"github.com/aminediro/iouring_server/server"
)

func main() {
	server.IncreaseResources()
	sigChannel := make(chan os.Signal, 1)

	l, err := server.MKRingListener(":8000")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening ring %#v\n", l)
	l.Listen()

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
	<-sigChannel

	l.Close()
	fmt.Println("Thanks for using Golang!")
}
