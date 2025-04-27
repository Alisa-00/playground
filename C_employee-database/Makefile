TARGET_SRV = bin/dbserver
TARGET_CLI = bin/dbclient
DB_FILE = mytestdb.db

SRC_SRV = $(wildcard src/database/*.c)
OBJ_SRV = $(SRC_SRV:src/database/%.c=obj/database/%.o)

SRC_CLI = $(wildcard src/client/*.c)
OBJ_CLI = $(SRC_CLI:src/client/%.c=obj/client/%.o)

init: clean default
	./$(TARGET_SRV) -f ./$(DB_FILE) -n

default: $(TARGET_SRV) $(TARGET_CLI)

testserver: init
	./$(TARGET_SRV) -f ./$(DB_FILE) -a "Alisa,Roer,100"
	./$(TARGET_SRV) -f ./$(DB_FILE) -a "Becky,Celdic,150"
	./$(TARGET_SRV) -f ./$(DB_FILE) -a "Claire,Saint Arkh,220"
	./$(TARGET_SRV) -f ./$(DB_FILE) -a "Dreichels,Heimdallr,280"
	./$(TARGET_SRV) -f ./$(DB_FILE) -a "Estelle,Rolent,330"
	./$(TARGET_SRV) -f ./$(DB_FILE) -a "Ferris,Ordis,350"
	./$(TARGET_SRV) -f ./$(DB_FILE) -p 5555 -a "Estelle,Rolent,500"

testclient: testclient1 testclient2 testclient3 testclient4

# add hours
testclient1:
	echo "Testing add hours"
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -l
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -s 1,23
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -s 2,10
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -s 7,55
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -l

# delete by id
testclient2:
	echo "Testing remove by id"
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -l
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -t 4
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -t 6
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -l

# delete by name
testclient3:
	echo "Testing remove by name"
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -l
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -r Estelle
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -l

# edit
testclient4:
	echo "Testing edit by id"
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -l
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -e 3,.,Heimdallr,.
	./$(TARGET_CLI) -h 127.0.0.1 -p 5555 -l

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