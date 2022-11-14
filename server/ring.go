package server

//#cgo LDFLAGS: -luring
//#include "server.h"
//#include <liburing.h>
import "C"
import "fmt"

type RINGOP int

const (
	ACCEPT RINGOP = 1 + iota
	READ
	WRITE
)

type IOURing struct {
	socketFD int
	conns    []int
}

func MkRing(socketFD int) (*IOURing, error) {
	err := int(C.ring_init())
	if err < 0 {
		return nil, fmt.Errorf("could not create the IOURING instance")
	}
	return &IOURing{
		socketFD: socketFD,
		conns:    make([]int, 10000),
	}, nil
}

func (r *IOURing) Run() {
	go r.accept()
	go r.getConn()
}

func (r *IOURing) accept() {
	C.ring_accept(C.int(r.socketFD))
}

func (r *IOURing) getConn() {
	// TODO : Hot loop do in C and callback into go ?
	for {
		connFD := int(C.completion_entry())
		r.conns = append(r.conns, connFD)
		if len(r.conns)%10 == 0 {
			fmt.Println("NB conns : ", len(r.conns))
		}
	}

}

func (r *IOURing) Close() {
	C.ring_close()
}
