#include <poll.h>
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netinet/in.h>

#include "db_poll.h"
#include "common.h"
#include "parse.h"

void init_clients(ClientState_t *clients) {
    int i=0;
    for (i = 0;i < MAX_CLIENTS; i++) {
        clients[i].fd = -1;
        clients[i].state = STATE_NEW;
        memset(&clients[i].buffer, '\0', BUFFER_SIZE);
    }
}

int find_free_slot(ClientState_t *ClientStates) {
    int i=0;
    for (i = 0;i < MAX_CLIENTS; i++) {
        if (ClientStates[i].fd == STATUS_ERROR) {
            return i;
        }
    }
    return STATUS_ERROR;
}

int find_slot_by_fd(int fd, ClientState_t *ClientStates) {
    int i=0;
    for (i = 0;i < MAX_CLIENTS; i++) {
        if (ClientStates[i].fd == fd) {
            return i;
        }
    }
    return STATUS_ERROR;
}

void fsm_reply_hello(ClientState_t *client, db_protocol_header_t *header) {
    header->type = htonl(MSG_HELLO_RESP);
    header->len = htons(1);
    db_protocol_hello_resp *hello = (db_protocol_hello_resp*)&header[1];
    hello->protocol = htons(PROTOCOL_VER);

    write(client->fd, header, sizeof(db_protocol_header_t) + sizeof(db_protocol_hello_resp));
}

void fsm_reply_err(ClientState_t *client, db_protocol_header_t *header) {
    header->type = htonl(MSG_ERROR);
    header->len = htons(0);

    write(client->fd, header, sizeof(db_protocol_header_t));
}

void fsm_reply_add(ClientState_t *client, db_protocol_header_t *header) {
    header->type = htonl(MSG_EMPLOYEE_ADD_RESP);
    header->len = htons(1);

    write(client->fd, header, sizeof(db_protocol_header_t));
}

void fsm_reply_list(ClientState_t *client, db_protocol_header_t *header, struct dbheader_t *database_header, struct employee_t *employees) {

    header->type = htonl(MSG_EMPLOYEE_LIST_RESP);
    header->len = htons(database_header->count);

    write(client->fd, header, sizeof(db_protocol_header_t));
    db_protocol_list_resp *employee = (db_protocol_list_resp*)&header[1];

    int i = 0;
    for (i=0; i<database_header->count; i++) {
        strncpy(employee->name, employees[i].name, sizeof(employees[0].name));
        strncpy(employee->address, employees[i].address, sizeof(employees[0].address));
        employee->hours = htonl(employees[i].hours);
        write(client->fd, &header[1], sizeof(db_protocol_list_resp));
    }
}

void fsm_reply_del(ClientState_t *client, db_protocol_header_t *header) {
    header->type = htonl(MSG_EMPLOYEE_DEL_RESP);
    header->len = htons(1);

    write(client->fd, header, sizeof(db_protocol_header_t));
}

int handle_client_fsm(struct dbheader_t *database_header, struct employee_t **employees, ClientState_t *client, char* filepath) {
    db_protocol_header_t *header = (db_protocol_header_t*)client->buffer;
    header->type = ntohl(header->type);
    header->len = ntohs(header->len);

    if (client->state == STATE_HELLO) {
        if (header->type != MSG_HELLO_REQ || header->len != 1) {
            printf("Didn't get MSG_HELLO in HELLO state\n");
            fsm_reply_err(client, header);
            return STATUS_ERROR;
        }

        db_protocol_hello_req* hello = (db_protocol_hello_req*)&header[1];
        hello->protocol = ntohs(hello->protocol);
        if (hello->protocol != PROTOCOL_VER) {
            printf("Protocol version mismatch\n");
            fsm_reply_err(client, header);
            return STATUS_ERROR;
        }

        fsm_reply_hello(client, header);
        client->state = STATE_MSG;
        printf("Client upgraded to STATE_MSG\n");
    }

    if (client->state == STATE_MSG) {

        if (header->type == MSG_EMPLOYEE_DEL_REQ) {
            db_protocol_del_req* employee = (db_protocol_del_req*)&header[1];
            printf("Removing employees with name: %s\n", employee->name);
            if (remove_employee(database_header, employees, (char*)employee->name) == STATUS_ERROR) {
                printf("Error removing employees!\n");
                fsm_reply_err(client, header);
                return STATUS_ERROR;
            }

            printf("Employees with name %s have been removed succesfully!\n", employee->name);
            fsm_reply_del(client, header);
            client->state = STATE_DISCONNECTED;
        }

        if (header->type == MSG_EMPLOYEE_ADD_REQ) {
            db_protocol_add_req* employee = (db_protocol_add_req*)&header[1];
            printf("Adding employee: %s\n", employee->data);

            if (add_employee(database_header, employees, (char*)employee->data) == STATUS_ERROR) {
                printf("Error adding new employee!\n");
                fsm_reply_err(client, header);
                return STATUS_ERROR;
            }

            printf("Employee was added succesfully!\n");
            fsm_reply_add(client, header);
            client->state = STATE_DISCONNECTED;
            output_file(database_header, *employees, filepath);

        }

        if (header->type == MSG_EMPLOYEE_LIST_REQ) {
            printf("Sending employee list..\n");
            fsm_reply_list(client, header, database_header, *employees);
            client->state = STATE_DISCONNECTED;
        }
    }

    return STATUS_SUCCESS;

}