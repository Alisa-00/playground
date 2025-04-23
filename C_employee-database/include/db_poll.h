#ifndef DB_POLL_H
#define DB_POLL_H

#include "parse.h"
#include "common.h"

#define BACKLOG 10
#define MAX_CLIENTS 256
#define BUFFER_SIZE 4096

typedef enum {
    STATE_NEW,
	STATE_CONNECTED,
	STATE_DISCONNECTED,
	STATE_HELLO,
	STATE_MSG
} State_enum;

typedef struct {
    int fd;
    State_enum state;
    char buffer[BUFFER_SIZE];
} ClientState_t;

void init_clients(ClientState_t *ClientStates);
int find_free_slot(ClientState_t *ClientStates);
int find_slot_by_fd(int fd, ClientState_t *ClientStates);
int handle_client_fsm(struct dbheader_t *header, struct employee_t **employees, ClientState_t *client, char *filepath);

#endif