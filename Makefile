.PHONY: build test run install tidy lint

MODULE := github.com/codebymaribel/eva-ai
BINARY := eva-ai
CMD    := ./cmd/eva-ai

# Build the binary to ./bin/
build:
	go build -o bin/$(BINARY) $(CMD)

# Run all tests with verbose output
test:
	go test ./... -v -count=1

# Run tests with coverage report
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run the CLI directly (pass ARGS="..." to forward flags)
run:
	go run $(CMD) $(ARGS)

# Install binary to $GOPATH/bin
install:
	go install $(CMD)

# Tidy modules
tidy:
	go mod tidy

# Dry-run example
dry-run:
	go run $(CMD) install --agent claude-code,cursor --dry-run

# Full install example (no changes, just for testing)
example:
	go run $(CMD) install --agent claude-code --preset full --dry-run