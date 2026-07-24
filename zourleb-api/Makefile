.PHONY: run worker build tidy test fmt docker-up docker-down

run:
	go run ./cmd/api

worker:
	go run ./cmd/worker

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down
