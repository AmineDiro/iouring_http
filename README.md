# C10M using IOURING

## Goal: 
@eranyanay used `epoll` + ulimit changes to push 1M websocket connections in GO.

We want to use  the new`io_uring` ring buffer to build highly concurrent ws server written in GO. 

For now we will use `liburing` which means we are using CGO. 

## system design : 

ringListener 
-> Http server will serve from the listener 
-> The listener will spawn a handler func and get the ringConnection

**OPTION 1 :**
-> The goroutine could read directly from the conn :
    We could maintain a map[FD]chan []bytes. The callbacks will push data to a chan per conn. Internally the RingConn would read from this channel
    This is not scalable => TOOO many goroutines and memory footprint

**OPTION 2:**
-> Handle the connection in a separate structure: (RingConnHandler???) 
-> Maintain a list of conns to read from. 
-> Submits read request to the ring via a channel + C call. 
-> Retrieve read completion results from a channel 


## TODO :
TCP based iouring server: 
- [x] Init ring in C
- [x] Ring based TCP listener 
- [ ] Stop all goroutines before closing the ring
- [ ] Read continuously from conns ??
Websocket : 
- [ ] Upgrade the connection to WS ?? 