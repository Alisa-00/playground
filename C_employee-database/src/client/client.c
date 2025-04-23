#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <poll.h>

#include "common.h"
#include "db_poll.h"

int send_hello(int socket) {
    char message_buffer[BUFFER_SIZE] = {0};

    db_protocol_header_t *header = (db_protocol_header_t*)message_buffer;
    header->type = MSG_HELLO_REQ;
    header->len = 1;

    //Add protocol version!
    db_protocol_hello_req *hello = (db_protocol_hello_req*)&header[1];
    hello->protocol = PROTOCOL_VER;

    header->type = htonl(header->type);
    header->len = htons(header->len);
    hello->protocol = htons(hello->protocol);

    // Send hello msg and read response
    write(socket, message_buffer, sizeof(db_protocol_header_t) + sizeof(db_protocol_hello_req));
    ssize_t bytes_read = read(socket, message_buffer, sizeof(message_buffer));

    // handle response
    if (bytes_read <= 0) {
        perror("read");
        return STATUS_ERROR;
    }

    header = (db_protocol_header_t*)message_buffer;
    header->type = ntohl(header->type);
    header->len = ntohs(header->len);

    if (header->type == MSG_ERROR) {
        printf("Protocol mismatch\n");
        return STATUS_ERROR;
    }

    if (header->type == MSG_HELLO_RESP) {
        db_protocol_hello_resp *hello_resp = (db_protocol_hello_resp*)&header[1];
        hello_resp->protocol = ntohs(hello->protocol);
        printf("Server connected, protocol v%d\n", hello_resp->protocol);
    }

    return STATUS_SUCCESS;
}

int send_add_req(int socket, char *employee_string) {
    char message_buffer[BUFFER_SIZE] = {0};

    db_protocol_header_t *header = (db_protocol_header_t*)message_buffer;
    header->type = MSG_EMPLOYEE_ADD_REQ;
    header->len = 1;

    db_protocol_add_req *employee = (db_protocol_add_req*)&header[1];
    strncpy(&employee->data[0], employee_string, sizeof(employee->data));

    header->type = htonl(header->type);
    header->len = htons(header->len);

    // Send add request and read response
    write(socket, message_buffer, sizeof(db_protocol_header_t) + sizeof(db_protocol_add_req));
    ssize_t bytes_read = read(socket, message_buffer, sizeof(message_buffer));

    // handle response
    if (bytes_read <= 0) {
        perror("read");
        return STATUS_ERROR;
    }

    header = (db_protocol_header_t*)message_buffer;
    header->type = ntohl(header->type);
    header->len = ntohs(header->len);

    if (header->type == MSG_ERROR) {
        printf("Error received, add request failed.\n");
        return STATUS_ERROR;
    }

    if (header->type == MSG_EMPLOYEE_ADD_RESP) {
        printf("Employee was added succesfully!\n");
    }

    return STATUS_SUCCESS;
}

int send_list_req(int socket) {
    char message_buffer[BUFFER_SIZE] = {0};

    db_protocol_header_t *header = (db_protocol_header_t*)message_buffer;
    header->type = MSG_EMPLOYEE_LIST_REQ;
    header->len = 1;

    header->type = htonl(header->type);
    header->len = htons(header->len);

    // Send add request and read response
    write(socket, message_buffer, sizeof(db_protocol_header_t));
    ssize_t bytes_read = read(socket, header, sizeof(db_protocol_header_t));

    // handle response
    if (bytes_read <= 0) {
        perror("read");
        return STATUS_ERROR;
    }

    header = (db_protocol_header_t*)message_buffer;
    header->type = ntohl(header->type);
    header->len = ntohs(header->len);

    if (header->type == MSG_ERROR) {
        printf("Error received, list request failed.\n");
        return STATUS_ERROR;
    }

    if (header->type == MSG_EMPLOYEE_LIST_RESP) {
        printf("Listing employees:\n");
        db_protocol_list_resp *employee = (db_protocol_list_resp*)&header[1];

        int i=0;
        for (i=0; i<header->len; i++) {
            read(socket, employee, sizeof(db_protocol_list_resp));
            employee->hours = ntohl(employee->hours);
            printf("%s, %s, %d\n", employee->name, employee->address, employee->hours);
        }

    }

    return STATUS_SUCCESS;
}

int send_del_req(int socket, char *employee_name) {
    char message_buffer[BUFFER_SIZE] = {0};

    db_protocol_header_t *header = (db_protocol_header_t*)message_buffer;
    header->type = MSG_EMPLOYEE_DEL_REQ;
    header->len = 1;

    db_protocol_del_req *employee = (db_protocol_del_req*)&header[1];
    strncpy(&employee->name[0], employee_name, sizeof(employee->name));

    header->type = htonl(header->type);
    header->len = htons(header->len);

    // Send add request and read response
    write(socket, message_buffer, sizeof(db_protocol_header_t) + sizeof(db_protocol_del_req));
    ssize_t bytes_read = read(socket, message_buffer, sizeof(message_buffer));

    // handle response
    if (bytes_read <= 0) {
        perror("read");
        return STATUS_ERROR;
    }

    header = (db_protocol_header_t*)message_buffer;
    header->type = ntohl(header->type);
    header->len = ntohs(header->len);

    if (header->type == MSG_ERROR) {
        printf("Error received, delete request failed.\n");
        return STATUS_ERROR;
    }

    if (header->type == MSG_EMPLOYEE_DEL_RESP) {
        printf("All employees with name %s have been deleted!\n", employee_name);
    }

    return STATUS_SUCCESS;
}

int main(int argc, char *argv[]) {

    char *addarg = NULL;
    char *delarg = NULL;
    char *portarg = NULL;
    char *hostarg = NULL;
    int list = 0;
    unsigned short port = 0;

    int c;
    while ((c = getopt(argc, argv, "p:h:a:d:l")) != -1) {
        switch(c) {
            case 'a':
                addarg = optarg;
                break;
            case 'd':
                delarg = optarg;
                break;
            case 'h':
                hostarg = optarg;
                break;
            case 'l':
                list = 1;
                break;
            case 'p':
                portarg = optarg;
                port = atoi(portarg);
                break;
            case '?':
                printf("Unknown option: -%c\n", c);
                break;
            default:
                return STATUS_ERROR;
        }
    }

    if (port == 0) {
        printf("Bad port: %s\n", portarg);
        return STATUS_ERROR;
    }

    if (hostarg == NULL) {
        printf("Must specify host with -h\n");
        return STATUS_ERROR;
    }

    struct sockaddr_in serverInfo = {0};
    serverInfo.sin_family = AF_INET;
    serverInfo.sin_addr.s_addr = inet_addr(hostarg);
    serverInfo.sin_port = htons(port);

    int server_socket = socket(AF_INET, SOCK_STREAM, 0);
    if (server_socket == STATUS_ERROR) {
        perror("socket");
        return STATUS_ERROR;
    }

    if (connect(server_socket, (struct sockaddr*)&serverInfo, sizeof(serverInfo)) == STATUS_ERROR) {
        perror("connect");
        close(server_socket);
        return STATUS_ERROR;
    }

    if (send_hello(server_socket) == STATUS_ERROR) {
        printf("Error establishing connection\n");
        close(server_socket);
        return STATUS_ERROR;
    }

    if (addarg != NULL) {
        if (send_add_req(server_socket, addarg) == STATUS_ERROR) {
            printf("Error trying to add new employee!\n");
            close(server_socket);
            return STATUS_ERROR;
        }
    }

    if (delarg != NULL) {
        if (send_del_req(server_socket, delarg) == STATUS_ERROR) {
            printf("Error trying to add new employee!\n");
            close(server_socket);
            return STATUS_ERROR;
        }
    }

    if (list > 0) {
        if (send_list_req(server_socket) == STATUS_ERROR) {
            printf("Error trying to add new employee!\n");
            close(server_socket);
            return STATUS_ERROR;
        }
    }

    return STATUS_SUCCESS;
}