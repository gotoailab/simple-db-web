# 数据库管理工具

一个使用 Go 和 Go Template 实现的现代化数据库管理 Web 工具，支持 MySQL 等数据库。

## 功能特性

- ✅ 建立数据库连接（支持 DSN 或表单输入）
- ✅ 列出数据库中的所有表
- ✅ 展示数据表的结构
- ✅ 支持数据表的查询（分页）
- ✅ 支持对数据表某一行数据的编辑
- ✅ 支持对数据表某一行数据的删除
- ✅ 支持写查询 SQL，展示查询结果
- ✅ 模块化设计，易于扩展支持其他数据库
- ✅ 现代化 UI，黄色主题（类似 Beekeeper Studio）

## 编译

```bash
go build -o dbweb
```

## 运行

```bash
./dbweb
```

服务器将在 `http://localhost:8080` 启动。

## 使用说明

1. **连接数据库**
   - 选择数据库类型（目前支持 MySQL）
   - 选择连接方式：
     - DSN 连接字符串：直接输入完整的 DSN，如 `user:password@tcp(host:port)/database`
     - 表单输入：分别输入主机、端口、用户名、密码、数据库名
   - 点击"连接"按钮

2. **查看表列表**
   - 连接成功后，左侧会显示所有数据表
   - 点击表名可以查看表的数据和结构

3. **查看表数据**
   - 在"数据"标签页中查看表的数据
   - 支持分页浏览
   - 可以编辑或删除某一行数据

4. **查看表结构**
   - 在"结构"标签页中查看表的 CREATE TABLE 语句

5. **执行 SQL 查询**
   - 在"SQL查询"标签页中输入 SQL 语句
   - 支持 SELECT、UPDATE、DELETE、INSERT 等操作
   - 点击"执行查询"按钮执行

## 项目结构

```
dbweb/
├── main.go              # 主程序入口
├── database/            # 数据库接口和实现
│   ├── interface.go     # Database 接口定义
│   └── mysql.go         # MySQL 实现
├── handlers/            # HTTP 处理器
│   └── handlers.go      # 路由和处理器
├── templates/           # HTML 模板
│   └── index.html       # 主页面
├── static/              # 静态资源
│   ├── style.css        # 样式文件
│   └── app.js           # 前端 JavaScript
└── README.md            # 说明文档
```

## 扩展支持其他数据库

要实现其他数据库的支持，只需：

1. 在 `database/` 目录下创建新的实现文件（如 `postgresql.go`）
2. 实现 `Database` 接口
3. 在 `handlers/handlers.go` 的 `Connect` 函数中添加新的数据库类型支持

示例：

```go
// database/postgresql.go
type PostgreSQL struct {
    db *sql.DB
}

func (p *PostgreSQL) Connect(dsn string) error {
    // 实现连接逻辑
}

// ... 实现其他接口方法
```

## 技术栈

- **后端**: Go 1.23+
- **数据库驱动**: github.com/go-sql-driver/mysql
- **前端**: 原生 JavaScript + CSS
- **模板**: Go Template

## 注意事项

- 目前仅支持 MySQL 数据库
- 编辑和删除操作基于主键（PRI）
- SQL 查询中的字符串值已做基本的转义处理，但建议在生产环境中使用参数化查询

## 许可证

MIT

