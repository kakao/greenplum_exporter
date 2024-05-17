BINARY_NAME=greenplum_exporter
GOBUILD=go build
GOCLEAN=go clean

all: build

build:
	@if [ ! -d bin/ ]; then mkdir bin/; fi;
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME) -v

clean:
	@$(GOCLEAN)
	@rm -f bin/$(BINARY_NAME)
.PHONY: clean

run: build
	@bin/$(BINARY_NAME)
.PHONY: run
