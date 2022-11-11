#include <liburing.h>
#include <stdio.h>

#define QUEUE_DEPTH 8096

char* hello_world()
{
   return "Hello, world!";
}

int queue_init(){
    struct io_uring ring;
    return io_uring_queue_init(QUEUE_DEPTH,&ring,0);
}
