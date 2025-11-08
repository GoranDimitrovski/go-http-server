.PHONY: build down test test-coverage test-coverage-html test-coverage-view clean run help

# Default target
help:
	@echo "Available targets:"
	@echo "  build              - Build Docker container"
	@echo "  down                - Stop and remove Docker containers"
	@echo "  test                - Run all tests"
	@echo "  test-coverage       - Run tests with coverage"
	@echo "  test-coverage-html  - Generate HTML coverage report"
	@echo "  test-coverage-view  - Generate and open HTML coverage report"
	@echo "  clean               - Clean build artifacts and coverage files"
	@echo "  run                 - Run the application locally"
	@echo "  docker-test         - Run tests in Docker container"

# Docker commands
build:
	@docker-compose up --build -d

down:
	@docker-compose down

docker-test:
	@docker-compose exec app go test -v ./...

# Test commands
test:
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

test-coverage-html:
	@echo "Generating HTML coverage report..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-view: test-coverage-html
	@echo "Opening coverage report in browser..."
	@open coverage.html || xdg-open coverage.html || echo "Please open coverage.html manually"

# Cleanup
clean:
	@echo "Cleaning build artifacts..."
	@rm -f app
	@rm -f coverage.out coverage.html
	@rm -f *.log
	@find . -name "*.test" -type f -delete
	@echo "Cleanup complete"

# Run locally
run:
	@go run main.go
