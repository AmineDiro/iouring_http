#include <assert.h>
#include <string.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <liburing.h>
#include <stdio.h>
#include <stdlib.h>
#include "_cgo_export.h"

#define NENTRIES 2048
#define READ_SIZE 512

#define EVENT_TYPE_ACCEPT 0
#define EVENT_TYPE_READ 1
#define EVENT_TYPE_WRITE 2

struct io_uring ring;

struct request
{
    int event_type;
    int iovec_count;
    int client_socket_fd;
    // Scattered *io_base pointers +size
    struct iovec iov[];
};

static void fatal_error(const char *syscall)
{
    perror(syscall);
    exit(1);
}

static void multishot_accept_entry(int socket_fd, struct sockaddr_in *client_addr, socklen_t *client_addr_len)
{
    struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);

    // multishot will fire a cqe everytime a conn comes in
    // NOTE: we don't care abour the client_addr and client_addr_len
    io_uring_prep_multishot_accept(sqe, socket_fd, (struct sockaddr *)client_addr,
                                   client_addr_len, 0);

    struct request *req = (struct request *)malloc(sizeof(*req));
    req->event_type = EVENT_TYPE_ACCEPT;

    io_uring_sqe_set_data(sqe, (void *)req);
    int ret = io_uring_submit(&ring);
    assert(ret != -1);
}

int read_entry(int client_socket_fd)
{
    struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);
    // TODO : use a single buffer !!
    struct request *req = malloc(sizeof(*req) + sizeof(struct iovec));
    req->iov[0].iov_base = malloc(READ_SIZE);
    req->iov[0].iov_len = READ_SIZE;
    req->event_type = EVENT_TYPE_READ;
    req->client_socket_fd = client_socket_fd;
    memset(req->iov[0].iov_base, 0, READ_SIZE);

    // Using readv instead of read uses iovec buffer
    // good for scattered io read and writes
    io_uring_prep_readv(sqe, client_socket_fd, req->iov, 1, 0);
    io_uring_sqe_set_data(sqe, req);
    io_uring_submit(&ring);
    return 0;
};

int completion_entry()
{
    struct io_uring_cqe *cqe;
    int wait_success = io_uring_wait_cqe(&ring, &cqe);
    if (wait_success == 0)
    {
        // The FD of the connections
        int conn_fd = cqe->res;
        if (conn_fd < 0)
        {
            return -1;
        }
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

    // Enables kernel thread polling => NO SYSCALL for submission
    // Needs root priviliges
    // params.flags |= IORING_SETUP_SQPOLL;

    return io_uring_queue_init_params(NENTRIES, &ring, &params);
}

void ring_loop(int socket_fd)
{
    struct io_uring_cqe *cqe;
    int nb_conns = 0;

    struct sockaddr_in client_addr;
    socklen_t client_addr_len = sizeof(client_addr);

    multishot_accept_entry(socket_fd, &client_addr, &client_addr_len);

    while (1)
    {
        int ret = io_uring_wait_cqe(&ring, &cqe);
        struct request *req = (struct request *)cqe->user_data;
        if (ret < 0){
            fatal_error("io_uring_wait_cqe");
        }
        if (cqe->res < 0)
        {
            fprintf(stderr, "Async request failed: %s for event: %d\n",
                    strerror(-cqe->res), req->event_type);
            exit(1);
        }
        switch (req->event_type)
        {

        case EVENT_TYPE_ACCEPT:
            nb_conns += 1;
            if (nb_conns % 100 == 0)
            {
                fprintf(stderr, "Nconns: %d\n", nb_conns);
            }

            accept_callback(cqe->res);
            // DEBug
            read_entry(cqe->res);
            break;

        case EVENT_TYPE_READ:
            // TODO: call the Read from GOLANG and track the data
            read_callback(req->client_socket_fd,(char *)req->iov[0].iov_base, READ_SIZE);
            read_entry(req->client_socket_fd);

            free(req->iov[0].iov_base);
            free(req);
            break;
        default:
            break;
        }

        io_uring_cqe_seen(&ring, cqe);
    }
}

void ring_close()
{
    io_uring_queue_exit(&ring);
}
