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

func (r *IOURing) accept() {
	go C.ring_accept(C.int(r.socketFD))
}
