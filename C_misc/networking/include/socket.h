#ifndef SOCKET_H
#define SOCKET_H

#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netinet/in.h>

#define PORT 5555
#define BACKLOG 5

typedef enum {
	PROTO_HELLO,
} proto_type_e;

typedef struct {
	proto_type_e type;
	unsigned int len;
} proto_hdr_t;

#endif