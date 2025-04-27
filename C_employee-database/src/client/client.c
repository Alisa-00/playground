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

int send_add_employee_req(int socket, char *employee_string) {
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
        printf("Error received, add employee request failed.\n");
        return STATUS_ERROR;
    }

    if (header->type == MSG_EMPLOYEE_ADD_RESP) {
        printf("Employee was added succesfully!\n");
    }

    return STATUS_SUCCESS;
}

int send_add_hrs_id_req(int socket, char *hrsstring) {
    char message_buffer[BUFFER_SIZE] = {0};

    db_protocol_header_t *header = (db_protocol_header_t*)message_buffer;
    header->type = MSG_EMPLOYEE_ADD_HRS_REQ;
    header->len = 1;

    db_protocol_add_hrs_req *employee = (db_protocol_add_hrs_req*)&header[1];
    strncpy(&employee->data[0], hrsstring, sizeof(employee->data));

    header->type = htonl(header->type);
    header->len = htons(header->len);

    // Send add request and read response
    write(socket, message_buffer, sizeof(db_protocol_header_t) + sizeof(db_protocol_add_hrs_req));
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
        printf("Error received, add hours request failed.\n");
        return STATUS_ERROR;
    }

    if (header->type == MSG_EMPLOYEE_ADD_HRS_RESP) {
        printf("Hours have been added succesfully:\n");
        //db_protocol_add_hrs_resp *employee = (db_protocol_add_hrs_resp*)&header[1];
        //printf("%d:\t%s, %s, %d\n", employee->id, employee->name, employee->address, employee->hours);
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
            employee->id = ntohl(employee->id);
            employee->hours = ntohl(employee->hours);
            printf("%d:\t%s, %s, %d\n", employee->id, employee->name, employee->address, employee->hours);
        }

    }

    return STATUS_SUCCESS;
}

int send_del_name_req(int socket, char *employee_name) {
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

int send_del_id_req(int socket, int id) {

    char message_buffer[BUFFER_SIZE] = {0};

    db_protocol_header_t *header = (db_protocol_header_t*)message_buffer;
    header->type = MSG_EMPLOYEE_DEL_ID_REQ;
    header->len = 1;

    db_protocol_del_id_req *employee = (db_protocol_del_id_req*)&header[1];
    employee->id = id;

    header->type = htonl(header->type);
    header->len = htons(header->len);
    employee->id = htonl(employee->id);

    // Send add request and read response
    write(socket, message_buffer, sizeof(db_protocol_header_t) + sizeof(db_protocol_del_id_req));
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
        printf("Employee with id %d has been deleted!\n", id);
    }

    return STATUS_SUCCESS;
}

int send_edit_id_req(int socket, int id) {
    return STATUS_SUCCESS;
}

void print_usage(char *argv[]) {
	printf("Usage: %s [-h HOST] [-p PORT]\n", argv[0]);
	printf("  -h  -  (required) host to connect to\n");
	printf("  -p  -  (required) port to connect to\n");
	printf("  -l  -  list employees\n");
	printf("  -t [id] -  remove employee by id\n");
	printf("  -r [name] -  remove employees by name\n");
	printf("  -s [name],[hours] - add hours to employee by id\n");
	printf("  -a [name],[address],[hours] -  add employee to the database\n");
	printf("  -e [id],[name],[address],[hours] - edit employee by id. blank fields will be left unchanged\n");
}

int main(int argc, char *argv[]) {

    char *addarg = NULL;
    char *hrsarg = NULL;
    char *delnamearg = NULL;
    char *delidarg = NULL;
    char *portarg = NULL;
    char *hostarg = NULL;
    int list = 0;
    unsigned short port = 0;
    unsigned int id = 0;

    int c;
    while ((c = getopt(argc, argv, "a:h:lp:r:s:t:")) != -1) {
        switch(c) {
            case 'a':
                addarg = optarg;
                break;
            case 'r':
                delnamearg = optarg;
                break;
            case 't':
                delidarg = optarg;
                id = atoi(delidarg);
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
            case 's':
                hrsarg = optarg;
                break;
            case '?':
                printf("Unknown option: -%c\n", c);
                break;
            default:
                return STATUS_ERROR;
        }
    }

    if (port == 0) {
        print_usage(argv);
        return STATUS_ERROR;
    }

    if (hostarg == NULL) {
        print_usage(argv);
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
        if (send_add_employee_req(server_socket, addarg) == STATUS_ERROR) {
            printf("Error with add new employee request!\n");
            close(server_socket);
            return STATUS_ERROR;
        }
    }

    if (hrsarg != NULL) {
        if (send_add_hrs_id_req(server_socket, hrsarg) == STATUS_ERROR) {
            printf("Error with add hours request!\n");
            close(server_socket);
            return STATUS_ERROR;
        }
    }

    if (delnamearg != NULL) {
        if (send_del_name_req(server_socket, delnamearg) == STATUS_ERROR) {
            printf("Error with delete employee by name request!\n");
            close(server_socket);
            return STATUS_ERROR;
        }
    }

    if (delidarg != NULL && id > 0) {
        if (send_del_id_req(server_socket, id) == STATUS_ERROR) {
            printf("Error with delete employee by id request!\n");
            close(server_socket);
            return STATUS_ERROR;
        }
    }

    if (list > 0) {
        if (send_list_req(server_socket) == STATUS_ERROR) {
            printf("Error with list employees request!\n");
            close(server_socket);
            return STATUS_ERROR;
        }
    }

    return STATUS_SUCCESS;
}