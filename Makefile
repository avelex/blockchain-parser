.PHONY: local-start
local-start:
	go run cmd/main.go --config config.yaml

.PHONY: build
build:
	docker compose build

.PHONY: start
start: build
	docker compose up -d

.PHONY: down
down:
	docker compose down -v

.PHONY: test
test:
	go test -race ./...