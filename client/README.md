# SimpleDBWeb Client

这是一个完整的 SimpleDBWeb 客户端版本，包含用户认证和管理功能。

## 功能特性

- ✅ 用户登录和认证
- ✅ Session 管理（使用 SQLite 存储）
- ✅ 用户管理（管理员可以创建、编辑、删除用户）
- ✅ 修改密码功能
- ✅ 退出登录功能
- ✅ 通过 `SetCustomScript` 注入用户管理 UI（不侵入核心代码）

## 技术栈

- **Web 框架**: Gin
- **数据库**: modernc.org/sqlite（用于存储用户和 session）
- **认证**: Session-based authentication
- **密码加密**: bcrypt

## 安装和运行

### 1. 安装依赖

```bash
cd client
go mod tidy
```

### 2. 运行应用

```bash
go run main.go [数据库文件路径] [端口]
```

示例：
```bash
# 使用默认配置（client.db, :8080）
go run main.go

# 指定数据库文件路径
go run main.go /path/to/database.db

# 指定数据库文件路径和端口
go run main.go /path/to/database.db 8080
```

### 3. 访问应用

打开浏览器访问：`http://localhost:8080`

默认管理员账户：
- 用户名: `admin`
- 密码: `admin123`

**重要**: 首次运行后请立即修改默认管理员密码！

## 项目结构

```
client/
├── main.go              # 主程序入口
├── db.go                # 数据库初始化和表结构
├── auth.go              # 用户认证和 session 管理
├── middleware.go        # 认证中间件
├── handlers.go          # HTTP 请求处理器
├── go.mod               # Go 模块定义
├── templates/           # HTML 模板
│   └── login.html       # 登录页面
└── static/              # 静态文件
    └── user-management.js  # 用户管理 UI（通过 embed 引入）
```

## API 接口

### 认证相关

- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 退出登录
- `GET /api/auth/current` - 获取当前用户信息
- `POST /api/auth/password` - 修改密码

### 用户管理（需要管理员权限）

- `GET /api/users` - 获取所有用户列表
- `POST /api/users` - 创建新用户
- `PUT /api/users/:id` - 更新用户信息
- `DELETE /api/users/:id` - 删除用户

## 数据库结构

### users 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| username | TEXT | 用户名（唯一） |
| password_hash | TEXT | 密码哈希值 |
| is_admin | INTEGER | 是否为管理员（0/1） |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### sessions 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| session_id | TEXT | 会话 ID（唯一） |
| user_id | INTEGER | 用户 ID（外键） |
| username | TEXT | 用户名 |
| created_at | DATETIME | 创建时间 |
| expires_at | DATETIME | 过期时间 |

## 安全说明

1. **密码加密**: 使用 bcrypt 进行密码哈希
2. **Session 过期**: Session 默认 24 小时过期
3. **Cookie 安全**: Session ID 存储在 HttpOnly Cookie 中
4. **权限控制**: 用户管理功能仅管理员可访问

## 注意事项

1. 首次运行会自动创建默认管理员账户（admin/admin123）
2. 如果数据库文件不存在，会自动创建
3. 过期 session 会定期自动清理（每小时一次）
4. 用户管理 UI 通过 `SetCustomScript` 注入，不会修改核心项目代码

