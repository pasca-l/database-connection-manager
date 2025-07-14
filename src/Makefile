FILE = dbcm

.PHONY: build
build:
	go build -o ../bin/$(FILE) ./main.go

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o ../bin/$(FILE) ./main.go
