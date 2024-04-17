tidy:
	@echo "Cleaning up..."
	@go mod tidy

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run
