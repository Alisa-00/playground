#ifndef PARSE_H
#define PARSE_H

#define HEADER_MAGIC 0x616C6973
#define HEADER_VERSION 1

struct dbheader_t {
    unsigned int magic;
    unsigned short version;
    unsigned short count;
    unsigned int filesize;
};

struct employee_t {
    char name[256];
    char address[256];
    unsigned int hours;
};

void output_file(struct dbheader_t *dbHeader, struct employee_t *dbEmployeeList, char *filename);
void list_employees(struct dbheader_t *dbHeader, struct employee_t *dbEmployeeList);
int create_db_header(int fileDescriptor, struct dbheader_t **headerOut);
int validate_db_header(int fileDescriptor, struct dbheader_t **headerOut);
int read_employees(int fileDescriptor, struct dbheader_t *dbHeader, struct employee_t **employeesOut);
int add_employee(struct dbheader_t *dbHeader, struct employee_t **employeesOut, char *addstring);
int remove_employee(struct dbheader_t *dbHeader, struct employee_t **employees, char *removeString);
int add_hours(struct dbheader_t *dbHeader, struct employee_t *employees, char *addString);


#endif