# Makefile for shelver

# Variables
BINARY_NAME=shelver
MAIN_PACKAGE=./cmd/shelver
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-trimpath

.PHONY: all build clean test test-verbose test-coverage install uninstall help deps tidy

# Default target
all: clean test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for development (no optimizations, faster builds)
build-dev:
	@echo "Building $(BINARY_NAME) for development..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Development build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for current directory (like go build)
build-local:
	@echo "Building $(BINARY_NAME) in current directory..."
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Local build completed: ./$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -cover ./...
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

# Tidy up go.mod
tidy:
	@echo "Tidying go.mod..."
	$(GOMOD) tidy

# Install the binary to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully"
	@echo "Run '$(BINARY_NAME) --help' to get started"

# Install to user's local bin (no sudo required)
install-user: build
	@echo "Installing $(BINARY_NAME) to ~/bin..."
	@mkdir -p ~/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) ~/bin/$(BINARY_NAME)
	@chmod +x ~/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed to ~/bin"
	@echo "Make sure ~/bin is in your PATH"
	@echo "Run '$(BINARY_NAME) --help' to get started"

# Uninstall the binary from system
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled successfully"

# Uninstall from user's local bin
uninstall-user:
	@echo "Uninstalling $(BINARY_NAME) from ~/bin..."
	@rm -f ~/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled from ~/bin"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	$(GOCLEAN)
	@echo "Clean completed"

# Run the application (for development)
run: build-dev
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run with example arguments
run-example: build-dev
	@echo "Running $(BINARY_NAME) with example (dry run on testdata)..."
	@$(BUILD_DIR)/$(BINARY_NAME) --dryrun testdata/*.wav testdata/*.jpg

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code (requires golint)
lint:
	@echo "Linting code..."
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "golint not found. Install with: go install golang.org/x/lint/golint@latest"; \
	fi

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Run all checks (format, vet, test)
check: fmt vet test

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux amd64
	@echo "Building for Linux amd64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	
	# Linux arm64
	@echo "Building for Linux arm64..."
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	
	# macOS amd64
	@echo "Building for macOS amd64..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	
	# macOS arm64 (Apple Silicon)
	@echo "Building for macOS arm64..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	
	# Windows amd64
	@echo "Building for Windows amd64..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	
	@echo "Multi-platform build completed"
	@ls -la $(BUILD_DIR)/

# Create a release package
release: build-all
	@echo "Creating release packages..."
	@cd $(BUILD_DIR) && \
	for binary in $(BINARY_NAME)-*; do \
		if [[ $$binary == *.exe ]]; then \
			zip "$${binary%.exe}.zip" "$$binary"; \
		else \
			tar -czf "$$binary.tar.gz" "$$binary"; \
		fi; \
	done
	@echo "Release packages created in $(BUILD_DIR)/"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary (optimized)"
	@echo "  build-dev     - Build for development (faster, no optimizations)"
	@echo "  build-local   - Build in current directory"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  test          - Run tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  install       - Install to system (/usr/local/bin, requires sudo)"
	@echo "  install-user  - Install to ~/bin (no sudo required)"
	@echo "  uninstall     - Uninstall from system"
	@echo "  uninstall-user- Uninstall from ~/bin"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Build and run the application"
	@echo "  run-example   - Run with example arguments"
	@echo "  deps          - Install dependencies"
	@echo "  tidy          - Tidy go.mod"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golint)"
	@echo "  vet           - Vet code"
	@echo "  check         - Run fmt, vet, and test"
	@echo "  release       - Create release packages"
	@echo "  help          - Show this help message"