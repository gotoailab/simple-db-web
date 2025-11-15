.PHONY: run run-auth build build-all build-linux build-windows build-darwin clean tidy help

# 应用配置
BINARY_NAME = simple-db-web-client
BUILD_DIR = ./build
MAIN_PATH = main.go

# 构建参数
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION = $(shell go version | awk '{print $$3}')

# 运行配置
PORT ?= 8080
DB_PATH ?= client.db
LOG_FILE ?= 
PREFIX ?= 
DEBUG ?= false
AUTH ?= false
OPEN ?= false
CONNECTIONS ?=

help:
	@echo "SimpleDBWeb Client Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make run              - Run application (default: port 8080, no auth)"
	@echo "  make run-auth         - Run application with authentication enabled"
	@echo "  make build            - Build application for current platform"
	@echo "  make build-all        - Build for all platforms (Linux, Windows, macOS)"
	@echo "  make build-linux      - Build Linux version (amd64, arm64)"
	@echo "  make build-windows    - Build Windows version (amd64, arm64)"
	@echo "  make build-darwin     - Build macOS version (amd64, arm64)"
	@echo "  make clean            - Clean build files and database"
	@echo "  make tidy             - Update dependencies"
	@echo "  make help             - Show this help message"
	@echo ""
	@echo "Build Variables (can be overridden):"
	@echo "  VERSION=              - Application version (default: git tag or 'dev')"
	@echo "  BUILD_TIME=            - Build time (default: current time)"
	@echo "  COMMIT=                - Git commit hash (default: git rev-parse)"
	@echo ""
	@echo "Run Variables (can be overridden):"
	@echo "  PORT=8080             - Server port (e.g., 8080 or :8080)"
	@echo "  DB_PATH=client.db      - Database file path (only used when auth is enabled)"
	@echo "  LOG_FILE=              - Log file path (empty means no file logging)"
	@echo "  PREFIX=                - Route prefix (e.g., /v1, /api)"
	@echo "  DEBUG=false            - Enable debug mode (true/false)"
	@echo "  AUTH=false             - Enable authentication (true/false)"
	@echo "  OPEN=false             - Auto-open browser (true/false)"
	@echo "  CONNECTIONS=            - Path to YAML file with preset connections"
	@echo ""
	@echo "Examples:"
	@echo "  make run"
	@echo "  make run-auth"
	@echo "  make run PORT=9000 PREFIX=/v1"
	@echo "  make build VERSION=1.0.0"
	@echo "  make build-all VERSION=1.0.0"
	@echo "  make build-linux VERSION=1.0.0"

run:
	@echo "Starting SimpleDBWeb Client..."
	@echo "Port: $(PORT)"
	@go run . \
		-port $(PORT) \
		$$([ "$(DEBUG)" = "true" ] && echo "-debug") \
		$$([ -n "$(LOG_FILE)" ] && echo "-log $(LOG_FILE)") \
		$$([ -n "$(PREFIX)" ] && echo "-prefix $(PREFIX)") \
		$$([ "$(OPEN)" = "true" ] && echo "-open") \
		$$([ -n "$(CONNECTIONS)" ] && echo "-connections $(CONNECTIONS)")

run-auth:
	@echo "Starting SimpleDBWeb Client with Authentication..."
	@echo "Port: $(PORT)"
	@echo "Database: $(DB_PATH)"
	@go run . \
		-port $(PORT) \
		-auth \
		-db $(DB_PATH) \
		$$([ "$(DEBUG)" = "true" ] && echo "-debug") \
		$$([ -n "$(LOG_FILE)" ] && echo "-log $(LOG_FILE)") \
		$$([ -n "$(PREFIX)" ] && echo "-prefix $(PREFIX)") \
		$$([ "$(OPEN)" = "true" ] && echo "-open") \
		$$([ -n "$(CONNECTIONS)" ] && echo "-connections $(CONNECTIONS)")

# 构建标志
LDFLAGS = -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.Commit=$(COMMIT)' \
	-X 'main.GoVersion=$(GO_VERSION)'

# 构建当前平台
build:
	@echo "Building SimpleDBWeb Client for current platform..."
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit: $(COMMIT)"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 构建 Linux 平台 (amd64, arm64)
build-linux:
	@echo "Building SimpleDBWeb Client for Linux..."
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit: $(COMMIT)"
	@mkdir -p $(BUILD_DIR)
	@echo "  -> Linux amd64 (static, no CGO)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -tags netgo -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "  -> Linux arm64 (static, no CGO)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -tags netgo -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-*"

# 构建 Windows 平台 (amd64, arm64)
build-windows:
	@echo "Building SimpleDBWeb Client for Windows..."
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit: $(COMMIT)"
	@mkdir -p $(BUILD_DIR)
	@echo "  -> Windows amd64 (static, no CGO)..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -tags netgo -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "  -> Windows arm64 (static, no CGO)..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -tags netgo -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-*.exe"

# 构建 macOS 平台 (amd64, arm64)
build-darwin:
	@echo "Building SimpleDBWeb Client for macOS..."
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit: $(COMMIT)"
	@mkdir -p $(BUILD_DIR)
	@echo "  -> macOS amd64 (static, no CGO)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -tags netgo -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "  -> macOS arm64 (static, no CGO)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -tags netgo -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-*"

# 构建所有平台
build-all: build-linux build-windows build-darwin
	@echo ""
	@echo "=========================================="
	@echo "All platforms build complete!"
	@echo "Build directory: $(BUILD_DIR)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "=========================================="
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)-* 2>/dev/null || echo "No build files found"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f client.db
	@echo "Clean complete"

tidy:
	@echo "Updating dependencies..."
	@go mod tidy
	@echo "Dependencies updated"

