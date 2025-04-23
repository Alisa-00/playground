#include "socket.h"
#include "../common.h"

int create_socket() {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd == STATUS_ERROR) {
        perror("socket");
        return STATUS_ERROR;
    }

    return fd;
}

int bind_socket(int socket, struct sockaddr_in *serverAddress) {
    if (bind(socket, (struct sockaddr *)serverAddress, sizeof(*serverAddress)) == STATUS_ERROR) {
        perror("bind");
        close(socket);
        return STATUS_ERROR;
    }

    return STATUS_SUCCESS;
}

int listen_socket(int socket) {
    if (listen(socket, BACKLOG) == STATUS_ERROR) {
        perror("listen");
        close(socket);
        return STATUS_ERROR;
    }

    return STATUS_SUCCESS;
}

void handle_client(int socket) {
    char packetBuffer[4096] = {0};
    proto_hdr_t *header = (proto_hdr_t*)packetBuffer;

    header->type = htonl(PROTO_HELLO);
    header->len = sizeof(int);

    int realLength = header->len;
    header->len = htons(header->len);

    int *data = (int*)&header[1];
    *data = htonl(1);
    if (write(socket, header, sizeof(proto_hdr_t) + realLength) == STATUS_ERROR) {
        perror("write");
        close(socket);
        return;
    }

    printf("Message sent to client, protocol v1\n");
    return;
}

int main() {

    int serverSocket = create_socket();

    struct sockaddr_in serverInfo = {0};
    struct sockaddr_in clientInfo = {0};
    socklen_t clientAddrLen = sizeof(clientInfo);

    serverInfo.sin_family = AF_INET;
    serverInfo.sin_addr.s_addr = INADDR_ANY;
    serverInfo.sin_port = htons(PORT);

    if (serverSocket == STATUS_ERROR) return STATUS_ERROR;

    printf("Our socket FD is %d\n", serverSocket);

    int bound = bind_socket(serverSocket, &serverInfo);
    if (bound == STATUS_ERROR) return STATUS_ERROR;

    int listening = listen_socket(serverSocket);
    if (listening == STATUS_ERROR) return STATUS_ERROR;

    

    int clientSocket = accept(serverSocket, (struct sockaddr *)&clientInfo, &clientAddrLen);

    if (clientSocket == STATUS_ERROR) {
        perror("accept");
        close(serverSocket);
        return STATUS_ERROR;
    }

    handle_client(clientSocket);

    close(serverSocket);
    close(clientSocket);
    return STATUS_SUCCESS;
}