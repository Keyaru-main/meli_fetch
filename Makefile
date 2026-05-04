# Melifetch Makefile

.PHONY: all build install clean test help

# Variables
BINARY_NAME=melifetch
MAIN_PATH=cmd/melifetch/main.go
INSTALL_PATH=/usr/local/bin
GO=go
GOFLAGS=-v

# Default target
all: build

# Build binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ Build complete: ./$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "✓ Multi-platform build complete"

# Install binary
install: build
	@echo "📦 Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@if [ -w $(INSTALL_PATH) ]; then \
		mv $(BINARY_NAME) $(INSTALL_PATH)/; \
		echo "✓ Installed to $(INSTALL_PATH)/$(BINARY_NAME)"; \
	else \
		echo "⚠️  Need sudo for installation..."; \
		sudo mv $(BINARY_NAME) $(INSTALL_PATH)/; \
		echo "✓ Installed to $(INSTALL_PATH)/$(BINARY_NAME)"; \
	fi

# Uninstall binary
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@if [ -w $(INSTALL_PATH) ]; then \
		rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@echo "✓ Uninstalled"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	$(GO) clean
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GO) test -v ./...
	@echo "✓ Tests complete"

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	$(GO) fmt ./...
	@echo "✓ Format complete"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed"; \
		echo "Install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin"; \
	fi

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✓ Dependencies updated"

# Run the binary
run:
	@echo "🚀 Running $(BINARY_NAME)..."
	$(GO) run $(MAIN_PATH)

# Development mode (build and run)
dev: build
	./$(BINARY_NAME)

# Show help
help:
	@echo "MeliFetch Makefile Commands:"
	@echo ""
	@echo "  make build          - Build binary"
	@echo "  make build-all      - Build for multiple platforms"
	@echo "  make install        - Install binary to $(INSTALL_PATH)"
	@echo "  make uninstall      - Uninstall binary"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Lint code"
	@echo "  make deps           - Download dependencies"
	@echo "  make run            - Run without building"
	@echo "  make dev            - Build and run"
	@echo "  make help           - Show this help"
	@echo ""
