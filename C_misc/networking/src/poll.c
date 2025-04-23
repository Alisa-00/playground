#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <poll.h>

#define MAX_CLIENTS 256
#define PORT 5555
#define BUFFER_SIZE 4096

typedef enum {
    STATE_NEW,
    STATE_CONNECTED,
    STATE_DISCONNECTED
} State_enum;

typedef struct {
    int fd;
    State_enum state;
    char buffer[BUFFER_SIZE];
} ClientState_t;

ClientState_t ClientStates[MAX_CLIENTS];

void init_clients() {
    int i=0;
    for (i = 0;i < MAX_CLIENTS; i++) {
        ClientStates[i].fd = -1;
        ClientStates[i].state = STATE_NEW;
        memset(&ClientStates[i].buffer, '\0', BUFFER_SIZE);
    }
}

int find_free_slot() {
    int i=0;
    for (i = 0;i < MAX_CLIENTS; i++) {
        if (ClientStates[i].fd == -1) {
            return i;
        }
    }
    return -1;
}

int find_slot_by_fd(int fd) {
    int i=0;
    for (i = 0;i < MAX_CLIENTS; i++) {
        if (ClientStates[i].fd == fd) {
            return i;
        }
    }
    return -1;
}

int main() {

    int listen_fd, conn_fd, freeSlot;
    struct sockaddr_in server_addr, client_addr;
    socklen_t client_len = sizeof(client_addr);

    struct pollfd fds[MAX_CLIENTS+1];
    int nfds = 1;
    int opt = 1;

    init_clients();
    
    listen_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_fd == -1) {
        perror("socket");
        return -1;
    }

    if (setsockopt(listen_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) == -1) {
        perror("setsockopt");
        close(listen_fd);
        return -1;
    }

    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(PORT);

    if (bind(listen_fd, (struct sockaddr*)&server_addr, sizeof(server_addr)) == -1) {
        perror("bind");
        close(listen_fd);
        return -1;
    }

    if (listen(listen_fd, 10) == -1) { // 10 backlong
        perror("listen");
        close(listen_fd);
        return -1;
    }

    printf("Server listening on port %d\n", PORT);

    memset(fds, 0, sizeof(fds));
    fds[0].fd = listen_fd;
    fds[0].events = POLLIN;
    nfds = 1;

    while (1) {
        
        int ii=0;
        int i=0;
        for (i = 0;i < MAX_CLIENTS; i++) {
            if (ClientStates[i].fd != -1) {
                fds[ii].fd = ClientStates[i].fd;
                fds[ii].events = POLLIN;
                ii++;
            }
        }

        int n_events = poll(fds, nfds, -1); // -1 is no timeout
        if (n_events == -1) {
            perror("poll");
            close(listen_fd);
            return -1;
        }

        if (fds[0].revents & POLLIN) {
            conn_fd = accept(listen_fd, (struct sockaddr*) &client_addr, &client_len);
            if (conn_fd == -1) {
                perror("accept");
                continue;
            }

            printf("New connection from %s:%d\n",inet_ntoa(client_addr.sin_addr), ntohs(client_addr.sin_port));

            freeSlot = find_free_slot();
            if (freeSlot == -1) {
                printf("Server full. Closing the connection\n");
                close(conn_fd);
            } else {
                ClientStates[freeSlot].fd = conn_fd;
                ClientStates[freeSlot].state = STATE_CONNECTED;
                nfds++;
                printf("Slot %d has fd %d\n", freeSlot, ClientStates[freeSlot].fd);
            }

            n_events--;

        }

        int j = 0;
        for (j = 0;j <= nfds && n_events > 0; j++) {
            if (fds[j].revents & POLLIN) {
                n_events--;

                int fd = fds[j].fd;
                int slot = find_slot_by_fd(fd);
                ssize_t bytes_read = read(fd, &ClientStates[slot].buffer, sizeof(&ClientStates[slot].buffer));

                if (bytes_read > 0) {
                    printf("Read data from client! %s\n", ClientStates[slot].buffer);
                } else {
                    perror("read");
                    close(fd);
                    if (slot == -1) {
                        printf("Tried closing fd that doesnt exist?\n");
                    } else {
                        ClientStates[slot].fd = -1;
                        ClientStates[slot].state = STATE_DISCONNECTED;
                        printf("Client disconnected!\n");
                        nfds--;
                    }
                }
            }
        }

    }

    return 0;
}