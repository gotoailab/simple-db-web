.PHONY: run build clean help

# 默认数据库文件路径
DB_PATH ?= client.db
# 默认端口
PORT ?= 8080

help:
	@echo "SimpleDBWeb Client Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make run          - 运行应用（使用默认配置：client.db, :8080）"
	@echo "  make run DB_PATH=/path/to/db.db PORT=8080  - 运行应用（指定数据库和端口）"
	@echo "  make build        - 编译应用"
	@echo "  make clean        - 清理编译文件和数据库"
	@echo "  make tidy         - 更新依赖"
	@echo "  make help         - 显示此帮助信息"
	@echo ""
	@echo "Examples:"
	@echo "  make run"
	@echo "  make run DB_PATH=./data/app.db PORT=9090"

run:
	@echo "Starting SimpleDBWeb Client..."
	@echo "Database: $(DB_PATH)"
	@echo "Port: $(PORT)"
	@go run . $(DB_PATH) $(PORT)

build:
	@echo "Building SimpleDBWeb Client..."
	@go build -o simple-db-web-client main.go
	@echo "Build complete: ./simple-db-web-client"

clean:
	@echo "Cleaning..."
	@rm -f simple-db-web-client
	@rm -f client.db
	@echo "Clean complete"

tidy:
	@echo "Updating dependencies..."
	@go mod tidy
	@echo "Dependencies updated"

