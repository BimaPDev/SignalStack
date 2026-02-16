.PHONY: dev-up dev-down migrate-up run-api run-worker build-api build-worker lint test

dev-up:
	docker compose up -d postgres
	@echo "Waiting for Postgres…"
	@sleep 2
	$(MAKE) migrate-up

dev-down:
	docker compose down -v

migrate-up:
	@for f in db/migrations/*.sql; do \
		echo "Applying $$f …"; \
		PGPASSWORD=pqpass psql -h localhost -U pq -d pulsequeue -f "$$f"; \
	done

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
