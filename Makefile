.PHONY: help build up down stop restart logs shell test fmt health

help:
	@echo "make build    Build and start the app"
	@echo "make up       Start the app"
	@echo "make down     Stop and remove the app"
	@echo "make stop     Stop the app"
	@echo "make restart  Restart the app"
	@echo "make logs     Follow logs"
	@echo "make shell    Open a shell in the app container"
	@echo "make test     go vet + go test -race in a container"
	@echo "make fmt      gofmt in a container"
	@echo "make health   Hit the health endpoint"

build:
	docker compose up --build -d

up:
	docker compose up -d

down:
	docker compose down

stop:
	docker compose stop

restart:
	docker compose restart

logs:
	docker compose logs -f

shell:
	docker compose exec app sh

GO = docker run --rm -v "$(CURDIR):/app" -w /app golang:1.22

test:
	$(GO) sh -c "go vet ./... && go test -race ./..."

fmt:
	$(GO) gofmt -l -w .

health:
	curl -s http://localhost:8000/health
