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
