#include "socket.h"
#include "../common.h"

void handle_server(int socket) {

	char packetBuffer[4096] = {0};
	proto_hdr_t *header = (proto_hdr_t*)packetBuffer;

    if (read(socket, header, sizeof(proto_hdr_t) + sizeof(int)) == STATUS_ERROR) {
        perror("read");
        close(socket);
        return;
    }

	header->type = ntohl(header->type); // unpack the type
	header->len = ntohs(header->len);

	int *data = (int*)&header[1];
	*data = ntohl(*data); // protocol version one, packed

    if (header->type != PROTO_HELLO) {
		printf("Protocol mismatch!\n");
		return;
	}

	if (*data != 1) {
		printf("Protocol version mismatch!\n");
		return;
	}

    printf("Successfully connected to the server, protocol v1\n");
	return;

}

int main(int argc, char *argv[]) {

    if (argc != 2) {
        printf("Usage: %s IP_ADDRESS\n", argv[0]);
        return STATUS_ERROR;
    }

    int serverSocket = socket(AF_INET, SOCK_STREAM, 0);
    if (serverSocket == STATUS_ERROR) {
        perror("socket");
        return STATUS_ERROR;
    }

    struct sockaddr_in serverInfo = {0};
    serverInfo.sin_family = AF_INET;
    serverInfo.sin_addr.s_addr = inet_addr(argv[1]);
    serverInfo.sin_port = htons(PORT);

    if (connect(serverSocket, (struct sockaddr *)&serverInfo, sizeof(serverInfo)) == STATUS_ERROR) {
        perror("connect");
        close(serverSocket);
        return STATUS_ERROR;
    }

    handle_server(serverSocket);
    close(serverSocket);

    return STATUS_SUCCESS;
}