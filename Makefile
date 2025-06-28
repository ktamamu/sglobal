.PHONY: build clean test install lint fmt vet

# Binary name
BINARY_NAME=sglobal

# Build the binary
build:
	go build -o $(BINARY_NAME) .

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	go test ./...

# Run linter (golangci-lint)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install linting tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all checks
check: fmt vet lint test

# Install to local bin
install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe .

# Create dist directory
dist:
	mkdir -p dist

# Release build
release: clean dist build-all