package server

//#cgo LDFLAGS: -luring
//#include "server.h"
//#include <liburing.h>
import "C"
import "net"

const (
	MAXCONN int = 2048
)

type RINGOP int

const (
	ACCEPT RINGOP = 1 + iota
	READ
	WRITE
)

type IOURing struct {
}

func MkRing() (int, error) {
	r := int(C.ring_init())
	return r, nil
}

type RingListener struct {
	ringFD int
	ring   *IOURing
	conns  chan net.Conn
}

func MKRingListener() (*RingListener, error) {
	return &RingListener{}, nil
}

// Called by the server to get the new accepted conns
func (rl *RingListener) Accept() (net.Conn, error) { return nil, nil }

// Closes the listener
func (rl *RingListener) Close() error { return nil }

// Gets the formated address
func (rl *RingListener) Addr() net.Addr { return nil }

// Start Listening by calling into the C liburing
// Either we continuously submit
func (rl *RingListener) Listen() {
	// get TCP Listener's FD
	// Send FD to
}
