package server

type RingConnHandler struct{}

func MkRingConnHandler() *RingConnHandler {
	return &RingConnHandler{}

}

// Starts the reading loop
// Iterate over the results
func (rh *RingConn) Start() {
}

// Adds the conn to the pool of conns to read from
func (rh *RingConn) Add(conn RingConn) {
}

// Removes the conn from the pool of conns
func (rh *RingConn) Remove(conn RingConn) {
}
