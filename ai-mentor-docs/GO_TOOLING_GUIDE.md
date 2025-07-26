# Go Tooling Guide for Proficient Development

## Essential Go Development Environment

This guide covers the essential tools, configurations, and workflows that will make you a proficient Go developer, especially for network programming and HTTP server development.

## Core Go Toolchain

### Go Installation and Version Management

#### Official Go Installation
```bash
# Download from https://golang.org/dl/
# Or use package managers:

# macOS (Homebrew)
brew install go

# Linux (apt)
sudo apt install golang-go

# Verify installation
go version
go env GOPATH
go env GOROOT
```

#### Version Management with g
```bash
# Install g (Go version manager)
curl -sSL https://git.io/g-install | sh -s

# Install and use specific Go versions
g install 1.24
g use 1.24
g list
```

### Go Workspace Setup

#### Modern Go Modules (Recommended)
```bash
# Initialize new module
go mod init your-module-name

# Add dependencies
go get github.com/stretchr/testify/assert
go get github.com/gorilla/mux

# Clean up dependencies
go mod tidy

# Verify dependencies
go mod verify

# View dependency graph
go mod graph
```

#### GOPATH vs Go Modules
```bash
# Check current settings
go env GOPATH
go env GOMODCACHE

# Modern development (Go 1.11+)
export GO111MODULE=on  # Default in Go 1.16+
```

## Essential Go Commands

### Building and Running

```bash
# Run directly
go run main.go
go run .
go run ./cmd/server

# Build binary
go build
go build -o server
go build -o bin/server ./cmd/server

# Install globally
go install
go install ./cmd/server

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o server-linux
GOOS=windows GOARCH=amd64 go build -o server.exe
GOOS=darwin GOARCH=arm64 go build -o server-mac-m1
```

### Testing

```bash
# Run tests
go test
go test ./...
go test -v
go test -short
go test -race

# Test with coverage
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Benchmark tests
go test -bench=.
go test -bench=BenchmarkHTTPServer -benchmem

# Test specific functions
go test -run TestServerStart
go test -run "TestServer.*"
```

### Code Quality and Analysis

```bash
# Format code
go fmt ./...
goimports -w .

# Vet for issues
go vet ./...

# Static analysis
golangci-lint run

# Documentation
go doc net/http
go doc -all net/http
godoc -http=:6060  # Local documentation server
```

## Advanced Go Tools

### Performance Analysis

#### pprof for Profiling
```bash
# CPU profiling
go build -o server
./server &
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profiling
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profiling
go tool pprof http://localhost:6060/debug/pprof/goroutine

# In your code, add pprof endpoint:
import _ "net/http/pprof"
```

#### Benchmarking and Profiling
```bash
# Run benchmarks with profiling
go test -bench=. -cpuprofile=cpu.prof
go test -bench=. -memprofile=mem.prof

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof

# Generate flame graphs
go tool pprof -http=:8080 cpu.prof
```

#### Execution Tracing
```bash
# Generate trace
go test -trace=trace.out
./server &
curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out

# Analyze trace
go tool trace trace.out
```

### Debugging Tools

#### Delve Debugger
```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug your program
dlv debug
dlv debug ./cmd/server

# Attach to running process
dlv attach <pid>

# Debug tests
dlv test
dlv test -- -test.run TestServerStart

# Remote debugging
dlv debug --headless --listen=:2345 --api-version=2
```

#### Race Detection
```bash
# Build with race detector
go build -race

# Run with race detection
go run -race main.go
go test -race ./...

# Environment variable
export GORACE="log_path=/tmp/race"
```

### Code Quality Tools

#### golangci-lint (Essential)
```bash
# Install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all linters
golangci-lint run

# Run specific linters
golangci-lint run --enable=gosec,gocyclo

# Configuration file (.golangci.yml)
```

#### Individual Linters
```bash
# Install useful linters
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Run individual tools
staticcheck ./...
gosec ./...
gocyclo -over 15 .
```

### Dependency Management

#### Security and Updates
```bash
# Check for known vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Update dependencies
go get -u
go get -u ./...
go get -u github.com/gorilla/mux

# Check outdated dependencies
go list -u -m all
```

#### Module Analysis
```bash
# Why is this dependency included?
go mod why github.com/some/dependency

# Download dependencies
go mod download

# Vendor dependencies (optional)
go mod vendor
go build -mod=vendor
```

## Development Environment Setup

### Editor/IDE Configuration

#### VS Code Setup
```json
// .vscode/settings.json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "go.formatTool": "goimports",
    "go.testFlags": ["-v", "-race"],
    "go.coverOnSave": true,
    "go.coverOnTestPackage": true,
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    }
}

// .vscode/launch.json for debugging
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/server",
            "env": {
                "PORT": "4221"
            }
        }
    ]
}
```

#### Essential VS Code Extensions
- Go (official)
- Go Test Explorer
- Go Outline
- REST Client (for testing HTTP servers)

#### Vim/Neovim Setup
```vim
" Install vim-go
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }

" Key mappings
au FileType go nmap <leader>r <Plug>(go-run)
au FileType go nmap <leader>b <Plug>(go-build)
au FileType go nmap <leader>t <Plug>(go-test)
au FileType go nmap <leader>c <Plug>(go-coverage)
```

### Shell Configuration

#### Useful Aliases
```bash
# Add to ~/.bashrc or ~/.zshrc
alias gt="go test -v"
alias gtr="go test -race -v"
alias gtc="go test -cover"
alias gb="go build"
alias gr="go run"
alias gf="go fmt ./..."
alias gv="go vet ./..."
alias gl="golangci-lint run"
alias gm="go mod tidy"
```

#### Environment Variables
```bash
# Add to shell profile
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org

# Development settings
export CGO_ENABLED=1
export GORACE="halt_on_error=1"
```

## HTTP Server Development Tools

### Testing HTTP Servers

#### httptest Package
```go
import (
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    
    handler(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

#### Load Testing Tools
```bash
# Install wrk
brew install wrk  # macOS
sudo apt install wrk  # Linux

# Basic load test
wrk -t12 -c400 -d30s http://localhost:4221/

# Install hey (Go-based)
go install github.com/rakyll/hey@latest
hey -n 1000 -c 10 http://localhost:4221/

# Apache Bench
ab -n 1000 -c 10 http://localhost:4221/
```

#### HTTP Client Testing
```bash
# Install httpie
pip install httpie

# Test requests
http GET localhost:4221/
http POST localhost:4221/echo/hello
http GET localhost:4221/user-agent User-Agent:curl/7.64.1

# Or use curl
curl -v localhost:4221/
curl -X POST localhost:4221/echo/hello
curl -H "User-Agent: custom" localhost:4221/user-agent
```

### Monitoring and Observability

#### Built-in Metrics
```go
import (
    "expvar"
    _ "net/http/pprof"
)

// Add metrics endpoint
http.Handle("/debug/vars", expvar.Handler())

// Custom metrics
var requests = expvar.NewInt("requests")
requests.Add(1)
```

#### Structured Logging
```bash
# Install popular logging libraries
go get go.uber.org/zap
go get github.com/sirupsen/logrus
go get log/slog  # Built-in since Go 1.21
```

## Project Structure and Tooling

### Standard Project Layout
```
your-project/
├── .golangci.yml
├── .gitignore
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── server/
│   └── config/
├── pkg/
│   └── http/
├── test/
├── scripts/
└── docs/
```

### Makefile for Common Tasks
```makefile
# Makefile
.PHONY: build test clean run lint

BINARY_NAME=server
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

test:
	go test -race -v ./...

test-cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

fmt:
	go fmt ./...
	goimports -w .

clean:
	go clean
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

deps:
	go mod download
	go mod verify

update-deps:
	go get -u ./...
	go mod tidy

bench:
	go test -bench=. -benchmem

profile-cpu:
	go test -bench=. -cpuprofile=cpu.prof
	go tool pprof cpu.prof

docker-build:
	docker build -t $(BINARY_NAME) .
```

### Configuration Files

#### .golangci.yml
```yaml
run:
  timeout: 5m
  modules-download-mode: readonly

linters-settings:
  gocyclo:
    min-complexity: 15
  govet:
    check-shadowing: true
  misspell:
    locale: US

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - errcheck
    - exhaustive
    - gochecknoinits
    - goconst
    - gocyclo
    - gofmt
    - goimports
    - golint
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - misspell
    - nakedret
    - rowserrcheck
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gosec
        - dupl
```

#### .gitignore for Go
```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out
coverage.html

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db

# Logs
*.log

# Local environment
.env
.env.local
```

## Continuous Integration

### GitHub Actions for Go
```yaml
# .github/workflows/go.yml
name: Go

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.24
        
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
          
    - name: Download dependencies
      run: go mod download
      
    - name: Run tests
      run: go test -race -coverprofile=coverage.out ./...
      
    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
        
    - name: Build
      run: go build -v ./...
      
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

## Development Workflow

### Daily Development Commands
```bash
# Start development session
go mod tidy
go test ./...
golangci-lint run

# During development
go run main.go
go test -run TestSpecificFunction
go test -v ./internal/handler

# Before committing
go fmt ./...
goimports -w .
go vet ./...
golangci-lint run
go test -race ./...

# Performance check
go test -bench=. -benchmem
```

### Debugging Workflow
```bash
# Add logging
go run main.go 2>&1 | tee debug.log

# Use delve for step debugging
dlv debug -- -addr :4221

# Profile running server
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Check for race conditions
go run -race main.go
```

## Advanced Tooling

### Custom Tools Development
```bash
# Create custom Go tools
go install -a std  # Rebuild standard library
go build -buildmode=plugin  # Build plugins

# Code generation
go generate ./...
go install golang.org/x/tools/cmd/stringer@latest
```

### Build Optimization
```bash
# Optimized builds
go build -ldflags="-s -w" -o server  # Strip debug info
go build -ldflags="-X main.version=1.0.0"  # Set variables

# Check binary size
ls -lh server
go tool nm server | wc -l  # Count symbols
```

This comprehensive tooling setup will make you highly productive in Go development, especially for HTTP server and network programming projects. Start with the core tools and gradually adopt the advanced ones as needed.