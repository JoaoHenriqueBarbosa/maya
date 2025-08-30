.PHONY: all build test clean dev wasm tinygo

# Variables
BINARY_NAME=maya
WASM_OUTPUT=dist/maya.wasm
GO_VERSION=1.24
GOFLAGS=-ldflags="-s -w"

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

all: test build

# Development server with hot reload
dev:
	@echo "$(GREEN)Starting development server...$(NC)"
	air -c .air.toml

# Build WASM with standard Go compiler
wasm:
	@echo "$(GREEN)Building WASM with Go $(GO_VERSION)...$(NC)"
	@mkdir -p dist
	GOOS=js GOARCH=wasm go build $(GOFLAGS) \
		-o $(WASM_OUTPUT) \
		./cmd/maya/...
	@echo "$(GREEN)WASM size: $$(du -h $(WASM_OUTPUT) | cut -f1)$(NC)"
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/

# Build optimized WASM with TinyGo
tinygo:
	@echo "$(YELLOW)Building optimized WASM with TinyGo...$(NC)"
	@mkdir -p dist
	tinygo build -o $(WASM_OUTPUT) \
		-target wasm \
		-no-debug \
		-size short \
		./cmd/maya/...
	@echo "$(GREEN)TinyGo WASM size: $$(du -h $(WASM_OUTPUT) | cut -f1)$(NC)"

# Build with WASI support (Go 1.24)
wasi:
	@echo "$(GREEN)Building with WASI 0.2 support...$(NC)"
	@mkdir -p dist
	GOOS=wasip1 GOARCH=wasm go build $(GOFLAGS) \
		-o dist/maya-wasi.wasm \
		./cmd/maya/...

# Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v -race -cover ./...

# Run WASM tests in browser
test-wasm:
	@echo "$(GREEN)Running WASM tests in browser...$(NC)"
	GOOS=js GOARCH=wasm go test -exec="$$(go env GOPATH)/bin/wasmbrowsertest" ./...

# Benchmarks with Go 1.24 testing.B.Loop
bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

# Clean build artifacts
clean:
	@echo "$(RED)Cleaning build artifacts...$(NC)"
	rm -rf dist/
	go clean -cache -testcache

# Install development tools
install-tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/agnivade/wasmbrowsertest@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

# Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...
	gofmt -s -w .

# Lint code
lint:
	@echo "$(GREEN)Linting code...$(NC)"
	golangci-lint run
	staticcheck ./...

# Generate WebGPU bindings
generate:
	@echo "$(GREEN)Generating WebGPU bindings...$(NC)"
	go generate ./...

# Serve example
serve: wasm
	@echo "$(GREEN)Serving example at http://localhost:8080$(NC)"
	cd examples/basic && python3 -m http.server 8080

# Help
help:
	@echo "Maya Framework - Build Commands"
	@echo ""
	@echo "  $(GREEN)make dev$(NC)        - Start development server with hot reload"
	@echo "  $(GREEN)make wasm$(NC)       - Build WASM with standard Go compiler"
	@echo "  $(GREEN)make tinygo$(NC)     - Build optimized WASM with TinyGo"
	@echo "  $(GREEN)make wasi$(NC)       - Build with WASI 0.2 support"
	@echo "  $(GREEN)make test$(NC)       - Run tests"
	@echo "  $(GREEN)make test-wasm$(NC)  - Run WASM tests in browser"
	@echo "  $(GREEN)make bench$(NC)      - Run benchmarks"
	@echo "  $(GREEN)make serve$(NC)      - Build and serve example"
	@echo "  $(GREEN)make clean$(NC)      - Clean build artifacts"
	@echo ""