#ifndef COMMON_H
#define COMMON_H

#define STATUS_ERROR -1
#define STATUS_SUCCESS 0
#define PROTOCOL_VER 100

#include <stdint.h>

typedef enum {
	MSG_HELLO_REQ,
    MSG_HELLO_RESP,
    MSG_EMPLOYEE_LIST_REQ,
    MSG_EMPLOYEE_LIST_RESP,
    MSG_EMPLOYEE_ADD_REQ,
    MSG_EMPLOYEE_ADD_RESP,
    MSG_EMPLOYEE_ADD_HRS_REQ,
    MSG_EMPLOYEE_ADD_HRS_RESP,
    MSG_EMPLOYEE_DEL_REQ,
    MSG_EMPLOYEE_DEL_RESP,
    MSG_EMPLOYEE_DEL_ID_REQ,
    MSG_EMPLOYEE_DEL_ID_RESP,
    MSG_EMPLOYEE_EDIT_REQ,
    MSG_EMPLOYEE_EDIT_RESP,
    MSG_ERROR
} db_protocol_type_enum;

typedef struct {
	db_protocol_type_enum type;
	uint16_t len;
} db_protocol_header_t;

typedef struct {
	uint16_t protocol;
} db_protocol_hello;

typedef struct {
	uint8_t data[1024];
} db_protocol_data_req;

typedef struct {
    uint32_t id;
	uint8_t name[256];
    uint8_t address[256];
    uint32_t hours;
} db_protocol_list_resp;

typedef struct {
	uint32_t id;
} db_protocol_del_id_req;

#endif