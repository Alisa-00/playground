#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <getopt.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <poll.h>

#include "common.h"
#include "file.h"
#include "parse.h"
#include "db_poll.h"

void print_usage(char *argv[]) {
	printf("Usage: %s [-n] [-f FILE] [-p PORT]\n", argv[0]);
	printf("  -n  -  create new database file\n");
	printf("  -f  -  (required) path to database file\n");
	printf("  -p  -  port to listen to. if absent server will run commands and exit\n");
	printf("  -l  -  list employees\n");
	printf("  -t [id] -  remove employee by id\n");
	printf("  -r [name] -  remove employees by name\n");
	printf("  -h [name],[hours] - add hours to employee by id\n");
	printf("  -a [name],[address],[hours] -  add employee to the database\n");
	printf("  -e [id],[name],[address],[hours] - edit employee by id. use '.' for any fields to be left unchanged\n");
}

void close_client(ClientState_t *client) {

	close(client->fd);
	client->fd = -1;
    client->state = STATE_DISCONNECTED;
    printf("Client disconnected!\n\n");

}

void poll_loop(unsigned short port, struct dbheader_t *database_header, struct employee_t *employees, char* filepath) {
	int listen_fd, conn_fd, freeSlot;
    struct sockaddr_in server_addr, client_addr;
    socklen_t client_len = sizeof(client_addr);
	ClientState_t ClientStates[MAX_CLIENTS] = {0};

    struct pollfd fds[MAX_CLIENTS+1];
    int nfds = 1;
    int opt = 1;

    init_clients(&ClientStates[0]);

    listen_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_fd == -1) {
        perror("socket");
        return;
    }

    if (setsockopt(listen_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) == STATUS_ERROR) {
        perror("setsockopt");
        close(listen_fd);
        return;
    }

    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(port);

    if (bind(listen_fd, (struct sockaddr*)&server_addr, sizeof(server_addr)) == STATUS_ERROR) {
        perror("bind");
        close(listen_fd);
        return;
    }

    if (listen(listen_fd, BACKLOG) == STATUS_ERROR) {
        perror("listen");
        close(listen_fd);
        return;
    }

    printf("Server listening on port %d\n", port);

    memset(fds, 0, sizeof(fds));
    nfds = 1;

    while (1) {
		fds[0].fd = listen_fd;
    	fds[0].events = POLLIN;

        int ii=1;
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
            return;
        }

        if (fds[0].revents & POLLIN) {
            conn_fd = accept(listen_fd, (struct sockaddr*) &client_addr, &client_len);
            if (conn_fd == -1) {
                perror("accept");
                continue;
            }

            printf("New connection from %s:%d\n",inet_ntoa(client_addr.sin_addr), ntohs(client_addr.sin_port));

            freeSlot = find_free_slot(&ClientStates[0]);
            if (freeSlot == -1) {
                printf("Server full. Closing the connection\n");
                close(conn_fd);
            } else {
                ClientStates[freeSlot].fd = conn_fd;
                ClientStates[freeSlot].state = STATE_HELLO;
                nfds++;
            }

            n_events--;

        }

        int j = 0;
        for (j = 0;j <= nfds && n_events > 0; j++) {
            if (fds[j].revents & POLLIN) {
                n_events--;

                int fd = fds[j].fd;
                int slot = find_slot_by_fd(fd, &ClientStates[0]);
                ssize_t bytes_read = read(fd, &ClientStates[slot].buffer, sizeof(ClientStates[slot].buffer));

                if (bytes_read <= 0) {
					printf("No new messages from client!\n");
					close_client(&ClientStates[slot]);
					nfds--;
					continue;
                }

				if (handle_client_fsm(database_header, &employees, &ClientStates[slot], filepath) == STATUS_ERROR) {
					printf("Error handling the message!\n");
					close_client(&ClientStates[slot]);
					nfds--;
					continue;
				}
            }
        }
    }
}

int main(int argc, char *argv[]) {

	char *filepath = NULL;
	char *addstring = NULL;
	char *removeString = NULL;
	char *addHours = NULL;
	char *portarg = NULL;
	bool newfile = false;
	bool listEmployees = false;
	int flag = 0;
	int port = 0;

	int dbFileDescriptor = -1;

	struct dbheader_t *dbHeader = NULL;
	struct employee_t *dbEmployeeList = NULL;

	while ((flag = getopt(argc, argv, "a:f:h:lnp:r:")) != -1) {

		switch(flag) {
			case 'a':
				addstring = optarg;
				break;
			case 'f':
				filepath = optarg;
				break;
			case 'h':
				addHours = optarg;
				break;
			case 'l':
				listEmployees = true;
				break;
			case 'n':
				newfile = true;
				break;
			case 'p':
				portarg = optarg;
				port = atoi(portarg);
				if (port == 0) {
					printf("bad port: %s\n",portarg);
				}
				break;
			case 'r':
				removeString = optarg;
				break;
			case '?':
				break;
            default:
                print_usage(argv);
                return STATUS_ERROR;
    	}
	}

	if (filepath == NULL) {
		print_usage(argv);
		return STATUS_ERROR;
	}

	if (newfile) {

		dbFileDescriptor = create_db_file(filepath);

		if (dbFileDescriptor == STATUS_ERROR) {
			printf("Error trying to create database file\n");
			return STATUS_ERROR;
		}

		if (create_db_header(dbFileDescriptor, &dbHeader) == STATUS_ERROR) {
			printf("Error trying to create database header\n");
			return STATUS_ERROR;
		}

	} else {

		dbFileDescriptor = open_db_file(filepath);

		if (dbFileDescriptor == STATUS_ERROR) {
			printf("Error trying to open database file\n");
			return STATUS_ERROR;
		}

		if (validate_db_header(dbFileDescriptor, &dbHeader) == STATUS_ERROR) {
			printf("Database header invalid!\n");
			return STATUS_ERROR;
		}

	}

	if (read_employees(dbFileDescriptor, dbHeader, &dbEmployeeList) == STATUS_ERROR) {
		printf("Error trying to read employees\n");
		return STATUS_ERROR;
	}

	//close until new output
	close(dbFileDescriptor);

	if (addstring != NULL) {

		if (dbEmployeeList == NULL) {
			perror("realloc");
			return STATUS_ERROR;
		}

		if (add_employee(dbHeader, &dbEmployeeList, addstring) == STATUS_ERROR) {
			printf("Error trying to add employee\n");
			return STATUS_ERROR;
		}

	}

	if (removeString != NULL) {
		if (remove_employee(dbHeader, &dbEmployeeList, removeString) == STATUS_ERROR) {
			printf("Error trying to remove employee\n");
			return STATUS_ERROR;
		}
	}

	if (addHours > 0) {
		if (add_hours(dbHeader, dbEmployeeList, addHours) == STATUS_ERROR) {
			printf("Error trying to add hours\n");
			return STATUS_ERROR;
		}
	}

	output_file(dbHeader, dbEmployeeList, filepath);

	if (listEmployees) {
		list_employees(dbHeader, dbEmployeeList);
	}

	if (port != 0) {
		poll_loop(port, dbHeader, dbEmployeeList, filepath);
	}

	free(dbHeader);
	free(dbEmployeeList);
	dbHeader = NULL;
	dbEmployeeList = NULL;

	return STATUS_SUCCESS;

}