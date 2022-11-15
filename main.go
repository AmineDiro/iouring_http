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
	defer l.Close()

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

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	// The "/" pattern matches everything, so we need to check
	// 	// that we're at the root here.
	// 	if req.URL.Path != "/" {
	// 		http.NotFound(w, req)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "hello io_uring!\n")
	// })

	// s := http.Server{Handler: mux}
	// if err := s.Serve(l); err != nil {
	// 	log.Fatal(err)
	// }
}
