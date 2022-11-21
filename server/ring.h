int ring_init();

void ring_loop(int fd);

int completion_entry();

int read_entry(int client_socket_fd);

// Clean  the ring
void ring_close();



