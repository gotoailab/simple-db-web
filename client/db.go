package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// DB 数据库连接
var DB *sql.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	// 创建表
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// 初始化默认管理员账户
	if err := initDefaultAdmin(); err != nil {
		return fmt.Errorf("failed to init default admin: %w", err)
	}

	return nil
}

// createTables 创建数据库表
func createTables() error {
	// 用户表
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		is_admin INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	// Session 表
	sessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT UNIQUE NOT NULL,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_sessions_session_id ON sessions(session_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
	`

	if _, err := DB.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	if _, err := DB.Exec(sessionsTable); err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	return nil
}

// initDefaultAdmin 初始化默认管理员账户
func initDefaultAdmin() error {
	// 检查是否已存在管理员
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = 1").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check admin users: %w", err)
	}

	// 如果已有管理员，不创建默认账户
	if count > 0 {
		return nil
	}

	// 创建默认管理员账户：admin/admin123
	hashedPassword, err := hashPassword("admin123")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = DB.Exec(
		"INSERT INTO users (username, password_hash, is_admin) VALUES (?, ?, ?)",
		"admin", hashedPassword, 1,
	)
	if err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

