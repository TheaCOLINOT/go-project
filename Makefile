.PHONY: build test db-up db-down

build:
	go build ./...

test:
	go test ./...

db-up:
	docker compose up -d db

db-down:
	docker compose down
