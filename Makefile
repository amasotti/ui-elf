# UI Elf - Makefile

# Binary name
BINARY_NAME=ui-elf

# Build directory
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Main package path
MAIN_PATH=cmd/ui-elf/main.go

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: ./$(BINARY_NAME)"

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows
	@echo "All platform builds complete"

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "Linux build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

# Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "macOS builds complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-*"

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Windows build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

# Install the binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(MAIN_PATH)
	@echo "Installation complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Run formatter
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Run vet
.PHONY: vet
vet:
	@echo "Running vet..."
	$(GOCMD) vet ./...

# Run all checks (fmt, vet, lint, test)
.PHONY: check
check: fmt vet lint test
	@echo "All checks passed"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f ui-elf-results.json
	@echo "Clean complete"

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Run the tool (example usage)
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME) with example parameters..."
	./$(BINARY_NAME) --component-type form --directory ./sample-files

# Display help
.PHONY: help
help:
	@echo "UI Elf - Makefile targets:"
	@echo ""
	@echo "  make build          - Build the binary for current platform"
	@echo "  make build-all      - Build binaries for all platforms (Linux, macOS, Windows)"
	@echo "  make build-linux    - Build binary for Linux"
	@echo "  make build-darwin   - Build binaries for macOS (amd64 and arm64)"
	@echo "  make build-windows  - Build binary for Windows"
	@echo "  make install        - Install binary to GOPATH/bin"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make clean          - Remove build artifacts and generated files"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make run            - Build and run with example parameters"
	@echo "  make help           - Display this help message"
	@echo ""

# Default target
.DEFAULT_GOAL := help
