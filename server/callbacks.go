package server

//#cgo LDFLAGS: -luring
//#include "ring.h"
//#include <liburing.h>
import "C"
import (
	"unsafe"
)

//export accept_callback
func accept_callback(connFD C.int) {
	acceptChan <- RingConn{clientSocketFD: int(connFD)}
}

//export read_callback
func read_callback(clientFD C.int, iovec *C.char, length C.int) {
	// TODO : Use some kind of preallocated buffer
	// The C.GoBytes create a bytes buffer and copies the iovec content
	// from memorypool
	// readLength := int(length)
	// buff <- MemPool
	// copy(buff, (*(*[129]byte)(unsafe.Pointer(iovec)))[:readLength:readLength])

	buff := C.GoBytes(unsafe.Pointer(iovec), length)
	connMap.mu.RLock()
	defer connMap.mu.RUnlock()
	connMap.chanMap[int(clientFD)] <- buff
}
