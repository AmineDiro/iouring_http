#include <netinet/in.h>
#include <arpa/inet.h>
#include <liburing.h>
#include <stdio.h>

#define QUEUE_DEPTH 8096

struct io_uring ring;

struct data
{
    __u8 opcode;
};

struct clientAddress
{
    struct sockaddr_in addr;
    unsigned int addr_len;
};

int ring_init()
{
    return io_uring_queue_init(QUEUE_DEPTH, &ring, 0);
}

void ring_accept(int fd, struct clientAddress *clientAddr)
{

    struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);

    while (true)
    {
        io_uring_prep_accept(sqe, fd, (struct sockaddr *)&clientAddr->addr,
                             &clientAddr->addr_len, 0);
        // TODO : Change this to hold connection data??
        struct data *d = calloc(sizeof(*d));
        d->opcode = 1;
        io_uring_sqe_set_data(sqe, d);
        io_uring_submit(&ring);
    }
}