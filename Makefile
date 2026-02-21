.PHONY: dev-up dev-down migrate-up run-api run-worker build-api build-worker lint test

dev-up:
	docker compose up -d --wait postgres
	$(MAKE) migrate-up

dev-down:
	docker compose down -v

migrate-up:
	docker compose exec postgres psql -U pq -d pulsequeue -f /docker-entrypoint-initdb.d/001_init.sql
	docker compose exec postgres psql -U pq -d pulsequeue -f /docker-entrypoint-initdb.d/002_indexes.sql

run-api:
	cd api && go run ./cmd/api

run-worker:
	cd worker && go run ./cmd/worker

build-api:
	cd api && go build -o bin/api ./cmd/api

build-worker:
	cd worker && go build -o bin/worker ./cmd/worker

lint:
	cd api && go vet ./...
	cd worker && go vet ./...

test:
	cd api && go test ./...
	cd worker && go test ./...

postgress: 
	docker run -d \
	  --name signalstack-db \
	  -e POSTGRES_USER=signalstack \
	  -e POSTGRES_PASSWORD=signalstack \
	  -e POSTGRES_DB=signalstack \
	  -p 5432:5432 \
	  postgres:17