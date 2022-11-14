# C10M using IOURING

## Goal: 
@eranyanay used `epoll` + ulimit changes to push 1M websocket connections in GO.

We want to use  the new`io_uring` ring buffer to build highly concurrent ws server written in GO. 

For now we will use `liburing` which means we are using CGO. 

## TODO :
TCP based iouring server: 
- [x] Init ring in C
- [x] Ring based TCP listener 
- [ ] Read continuously from conns
- [ ] Use the 
Websocket : 
- [ ] Upgrade the connection to WS ?? 