TARGET_SRV = bin/dbserver
TARGET_CLI = bin/dbclient

SRC_SRV = $(wildcard src/database/*.c)
OBJ_SRV = $(SRC_SRV:src/database/%.c=obj/database/%.o)

SRC_CLI = $(wildcard src/client/*.c)
OBJ_CLI = $(SRC_CLI:src/client/%.c=obj/client/%.o)

init: clean default

default: $(TARGET_SRV) $(TARGET_CLI)

run: clean default
	./$(TARGET_SRV) -f ./mynewdb.db -n -p 5555 &
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -a "Alice,Nameless Street 120,100"

add:
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -a "Boris,Russian Street 420,200"

kill:
	kill -9 $$(pidof dbserver)

clean:
	rm -f obj/database/*.o
	rm -f bin/*
	rm -f *.db

$(TARGET_SRV): $(OBJ_SRV)
	gcc -o $@ $?

$(OBJ_SRV): obj/database/%.o: src/database/%.c
	gcc -c $< -o $@ -Iinclude

$(TARGET_CLI): $(OBJ_CLI)
	gcc -o $@ $?

$(OBJ_CLI): obj/client/%.o: src/client/%.c
	gcc -c $< -o $@ -Iinclude