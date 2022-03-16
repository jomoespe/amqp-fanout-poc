EXECUTABLES = go
CHECK_EXECUTABLES := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error No $(exec) in PATH. Check README.md for build requirements)))

all: bin/produce bin/consume

clean:
	@ rm -fr bin/

bin/produce:
	@ go build -o bin/produce producer/main.go

bin/consume:
	@ go build -o bin/consume consumer/main.go

