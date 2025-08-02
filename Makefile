MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: build
build:
	$(MAKE) -C src build

.PHONY: start
start:
	$(MAKE) -C src build-linux
	docker compose up --build

.PHONY: client
client:
	docker compose exec db_client bash

.PHONY: lint
lint:
	docker run -t --rm -v $(MAKEFILE_DIR):/app -w /app/src golangci/golangci-lint:v2.3.0 golangci-lint run
