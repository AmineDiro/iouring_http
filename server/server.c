#include <string.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <liburing.h>
#include <stdio.h>
#include <stdlib.h>

#define NENTRIES 1024

struct io_uring ring;
struct io_uring_cqe *cqe;

void accept_entry(int socket_fd)
{
    struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);
    struct sockaddr_in client_addr;
    socklen_t client_addr_len = sizeof(client_addr);

    // multishot will fire a cqe everytime a conn comes in
    io_uring_prep_multishot_accept(sqe, socket_fd, (struct sockaddr *)&client_addr,
                                   &client_addr_len, 0);
    // TODO : Change this to hold connection data??
    int opCode = 1;
    io_uring_sqe_set_data(sqe, (void *)&opCode);
    io_uring_submit(&ring);
}

int completion_entry()
{
    int wait_success= io_uring_wait_cqe(&ring, &cqe);
    if (wait_success== 0)
    {
        // The FD of the connections
        int conn_fd = cqe->res;
        io_uring_cqe_seen(&ring, cqe);
        return conn_fd;
    }
    io_uring_cqe_seen(&ring, cqe);
    return -1;

}

int ring_init()
{
    struct io_uring_params params;
    memset(&params, 0, sizeof(params));
    // enables kernel thread polling => NO SYSCALL
    // Needs root priviliges
    // params.flags |= IORING_SETUP_SQPOLL;

    return io_uring_queue_init_params(NENTRIES, &ring, &params);
}

void ring_accept(int socket_fd)
{
    // Accept entries
    accept_entry(socket_fd);
}

void ring_close()
{
    io_uring_queue_exit(&ring);
}
