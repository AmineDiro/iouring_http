#include <string.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <liburing.h>
#include <stdio.h>
#include <stdlib.h>

#define NENTRIES 8096

struct io_uring ring;

void accept_entry(int socket_fd)
{
    struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);
    struct sockaddr_in client_addr;
    socklen_t client_addr_len = sizeof(client_addr);

    io_uring_prep_accept(sqe, socket_fd, (struct sockaddr *)&client_addr,
                         &client_addr_len, 0);
    // TODO : Change this to hold connection data??
    int opCode = 1;
    io_uring_sqe_set_data(sqe, (void *)&opCode);
    io_uring_submit(&ring);
}

int ring_init()
{
    struct io_uring_params params;
    memset(&params, 0, sizeof(params));
    // enables kernel thread polling => NO SYSCALL
    // params.flags |= IORING_SETUP_SQPOLL;

    return io_uring_queue_init_params(NENTRIES, &ring, &params);
}

int ring_accept(int socket_fd)
{
    fprintf(stderr, "Entered ring_accept");
    struct io_uring_cqe *cqe;

    // Accept entries
    while (true)
    {
        accept_entry(socket_fd);
        int peek = io_uring_peek_cqe(&ring, &cqe);
        if (!peek)
        {
        }
        io_uring_cqe_seen(&ring, cqe);
    }
}