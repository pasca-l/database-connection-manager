MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
FILE = dbcm

.PHONY: build
build:
	go build -o $(MAKEFILE_DIR)../bin/$(FILE) ./main.go

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(MAKEFILE_DIR)../bin/$(FILE) ./main.go
