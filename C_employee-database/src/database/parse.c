#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

#include <string.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <arpa/inet.h>
#include <fcntl.h>

#include "parse.h"
#include "common.h"

#define TEMP_DB_FILE "TEMP_DB_FILE.db"

void output_file(struct dbheader_t *dbHeader, struct employee_t *dbEmployeeList, char* filename) {
    // new file
    int tempFileDescriptor = open(TEMP_DB_FILE, O_RDWR | O_CREAT, 0644);

    unsigned short dbHeaderCount = dbHeader->count;

    //make copies
    struct dbheader_t *db_header_copy = malloc(sizeof(struct dbheader_t));
    memcpy(db_header_copy, dbHeader, sizeof(struct dbheader_t));
    struct employee_t *employees_copy = malloc(sizeof(struct employee_t) * dbHeaderCount);
    if (employees_copy == NULL) {
        perror("malloc");
        return;
    }
    memcpy(employees_copy, dbEmployeeList, sizeof(struct employee_t) * dbHeaderCount);



    // unpack header and write into file
    db_header_copy->magic = htonl(db_header_copy->magic);
    db_header_copy->version = htons(db_header_copy->version);
    db_header_copy->count = htons(db_header_copy->count);
    db_header_copy->filesize = htonl(db_header_copy->filesize);

    if (lseek(tempFileDescriptor, 0, SEEK_SET) == STATUS_ERROR) {
        perror("lseek");
        return;
    }

    if (write(tempFileDescriptor, db_header_copy, sizeof(struct dbheader_t)) == STATUS_ERROR) {
        perror("write");
        return;
    }

    if (lseek(tempFileDescriptor, sizeof(struct dbheader_t), SEEK_SET) == STATUS_ERROR) {
        perror("lseek");
        return;
    }

    // unpack employee data, write and repack for later use
    int i=0;
    for (i=0;i<dbHeaderCount;i++) {
        employees_copy[i].hours = htonl(employees_copy[i].hours);

        lseek(tempFileDescriptor, sizeof(struct dbheader_t) + sizeof(struct employee_t) * i, SEEK_SET);
        if (write(tempFileDescriptor, &employees_copy[i], sizeof(struct employee_t)) == STATUS_ERROR) {
            perror("write");
            return;
        }
    }

    free(db_header_copy);
    free(employees_copy);

    // delete db file, rename temp file to db filename
    // avoids leftover data at the end when the output file is smaller than the current db file
    close(tempFileDescriptor);
    remove(filename);
    rename(TEMP_DB_FILE, filename);

}

void list_employees(struct dbheader_t *dbHeader, struct employee_t *dbEmployeeList) {

    int i=0;
    for (i=0;i<dbHeader->count;i++) {
        printf("Employee %d:\n\tName: %s\n\tAddress: %s\n\tHours: %d\n\n", i+1, dbEmployeeList[i].name, dbEmployeeList[i].address, dbEmployeeList[i].hours);
    }
}

int create_db_header(int fileDescriptor, struct dbheader_t **headerOut) {

    struct dbheader_t *header = calloc(1, sizeof(struct dbheader_t));

    if (header == NULL) {

    }

    header->magic = HEADER_MAGIC;
    header->version = HEADER_VERSION;
    header->count = 0;
    header->filesize = sizeof(struct dbheader_t);

    *headerOut = header;

    return STATUS_SUCCESS;
}

int validate_db_header(int fileDescriptor, struct dbheader_t **headerOut) {

    if (fileDescriptor == STATUS_ERROR) {
        printf("Got invalid file descriptor!\n");
        return STATUS_ERROR;
    }

    struct dbheader_t *header = calloc(1, sizeof(struct dbheader_t));

    if (header == NULL) {
        perror("calloc");
        return STATUS_ERROR;
    }

    if (read(fileDescriptor, header, sizeof(struct dbheader_t)) == STATUS_ERROR) {
        perror("read");
        free(header);
        return STATUS_ERROR;
    }

    header->magic = ntohl(header->magic);
    header->version = ntohs(header->version);
    header->count = ntohs(header->count);
    header->filesize = ntohl(header->filesize);

    if (header->magic != HEADER_MAGIC) {
        printf("Got invalid magic number!\n");
        free(header);
        return STATUS_ERROR;
    }

    if (header->version != HEADER_VERSION) {
        printf("Got invalid version number!\n");
        free(header);
        return STATUS_ERROR;
    }

    struct stat dbstat = {0};
    fstat(fileDescriptor, &dbstat);

    if (header->filesize != dbstat.st_size) {
        printf("Corrupted database!\n");
        printf("Expected size: %d\nActual size: %d\n", header->filesize, dbstat.st_size);
        free(header);
        return STATUS_ERROR;
    }

    *headerOut = header;
    return STATUS_SUCCESS;
}

int read_employees(int fileDescriptor, struct dbheader_t *dbHeader, struct employee_t **employeesOut) {

    if (fileDescriptor == STATUS_ERROR) {
        printf("Got invalid file descriptor!\n");
        return STATUS_ERROR;
    }

    int count = dbHeader->count;
    struct employee_t *employees = calloc(count, sizeof(struct employee_t));

    if (employees == NULL) {
        perror("calloc");
        return STATUS_ERROR;
    }

    if (lseek(fileDescriptor, sizeof(struct dbheader_t), SEEK_SET) == STATUS_ERROR) {
        perror("lseek");
        return STATUS_ERROR;
    }

    if (read(fileDescriptor, employees, sizeof(struct employee_t) * count) == STATUS_ERROR) {
        perror("read");
        return STATUS_ERROR;
    }

    int i=0;
    for (i=0;i<count;i++) {
        employees[i].hours = ntohl(employees[i].hours);
    }

    *employeesOut = employees;

    return STATUS_SUCCESS;
}

int add_employee(struct dbheader_t *dbHeader, struct employee_t **employees_pointer, char *addstring) {

    char *employeeName = strtok(addstring, ",");
    if (employeeName == NULL) {
        printf("Wrong string format!\n");
        return STATUS_ERROR;
    }

    char *employeeAddress = strtok(NULL, ",");
    if (employeeAddress == NULL) {
        printf("Wrong string format!\n");
        return STATUS_ERROR;
    }

    char *employeeHours = strtok(NULL, ",");
    if (employeeHours == NULL) {
        printf("Wrong string format!\n");
        return STATUS_ERROR;
    }

    dbHeader->count = dbHeader->count+1;
	dbHeader->filesize = (dbHeader->filesize) + sizeof(struct employee_t);

    struct employee_t *employees = calloc(dbHeader->count, sizeof(struct employee_t));
    memcpy(employees, *employees_pointer, (dbHeader->count-1) * sizeof(struct employee_t));

    strncpy(employees[dbHeader->count-1].name, employeeName, sizeof(employees[dbHeader->count-1].name));
    strncpy(employees[dbHeader->count-1].address, employeeAddress, sizeof(employees[dbHeader->count-1].address));
    employees[dbHeader->count-1].hours = (unsigned int)atoi(employeeHours);

    *employees_pointer = employees;

    return STATUS_SUCCESS;
}

int add_hours(struct dbheader_t *dbHeader, struct employee_t *employees, char *addString) {

    char *employeeName = strtok(addString, ",");
    int employeeHours = atoi(strtok(NULL, ","));

    int i=0;
    for (i=0;i<dbHeader->count;i++) {
        if (strcmp(employees[i].name, employeeName) == 0) {
            employees[i].hours += employeeHours;
            return STATUS_SUCCESS;
        }
    }

    return STATUS_ERROR;
}

int remove_employee(struct dbheader_t *dbHeader, struct employee_t **employees, char *employeeName) {

    unsigned short count = dbHeader->count;
    struct employee_t *employeeList = *employees;
    bool removed = false;

    int i=0;
    for (i=0;i<count;i++) {
        if (strcmp(employeeList[i].name, employeeName) == 0) {

            for (int j = i; j < count - 1; ++j) {
                employeeList[j] = employeeList[j + 1];  // Shift employees left

                // decrease i once in case of consecutive removals
                if (j == count - 2) {
                    i--;
                }
            }

            count--;
            removed = true;
        }
    }

    if (!removed) {
        printf("No employee with name %s found!\n", employeeName);
        return STATUS_ERROR;
    }

    struct employee_t *newEmployeeList = realloc(employeeList, count * sizeof(struct employee_t));
    if (newEmployeeList == NULL && count > 0) {
        return STATUS_ERROR;
    }

    *employees = newEmployeeList;

    dbHeader->count = count;
    dbHeader->filesize = sizeof(struct dbheader_t) + count * sizeof(struct employee_t);
    return STATUS_SUCCESS;
}