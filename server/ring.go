package server

//#cgo LDFLAGS: -luring
//#include "server.h"
//#include <liburing.h>
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

type RINGOP int

const (
	ACCEPT RINGOP = 1 + iota
	READ
	WRITE
)

type IOURing struct {
	socketFD int
	nbconns  int
}

func MkRing(socketFD int) (*IOURing, error) {
	err := int(C.ring_init())
	if err < 0 {
		return nil, fmt.Errorf("could not create the IOURING instance")
	}
	return &IOURing{
		socketFD: socketFD,
	}, nil
}

func (r *IOURing) loop() {
	// Starts the ring loop
	C.ring_loop(C.int(r.socketFD))
}

func charToBytes(src *C.char, sz int) []byte {
	return C.GoBytes(unsafe.Pointer(src), C.int(sz))
}

//export Read_callback
func Read_callback(iovec *C.char, length C.int) {
	// readLength := int(length)
	// buff := make([]byte, readLength)

	// copy(buff, (*(*[1024]byte)(unsafe.Pointer(iovec)))[:readLength:readLength])
	buff := charToBytes(iovec, int(length))
	fmt.Printf("Received :%s\n", buff)

}

func (r *IOURing) getConn() {
	// TODO : Hot loop do in C and callback into go ?
	for {
		connFD := int(C.completion_entry())
		if connFD == -1 {
			log.Println("bad response")
		} else {
			r.nbconns += 1
			if r.nbconns%100 == 0 {
				log.Println("NB conns : ", r.nbconns)
			}
		}
	}

}

func (r *IOURing) Close() {
	C.ring_close()
}
