run: build
	@./bin/ginknown

build:
	@go build -o bin/ginknown cmd/main.go

test:
	@go test -v ./...

all: run

.PHONY: run build test all help

help:
	@echo "run : ./bin/ginknown"
	@echo "build : go build -o bin/ginknown cmd/main.go"
	@echo "test : go test -v ./..."
	@echo "all : run"