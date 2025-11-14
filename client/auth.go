package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Session 会话模型
type Session struct {
	ID        int       `json:"id"`
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// hashPassword 加密密码
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// checkPassword 验证密码
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// generateSessionID 生成会话ID
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// AuthenticateUser 验证用户登录
func AuthenticateUser(username, password string) (*User, error) {
	var user User
	var passwordHash string

	err := DB.QueryRow(
		"SELECT id, username, password_hash, is_admin, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &passwordHash, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid username or password")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if !checkPassword(passwordHash, password) {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}

// CreateSession 创建会话
func CreateSession(userID int, username string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour) // 24小时过期

	result, err := DB.Exec(
		"INSERT INTO sessions (session_id, user_id, username, expires_at) VALUES (?, ?, ?, ?)",
		sessionID, userID, username, expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get session ID: %w", err)
	}

	return &Session{
		ID:        int(id),
		SessionID: sessionID,
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}, nil
}

// GetSession 获取会话
func GetSession(sessionID string) (*Session, error) {
	var session Session
	err := DB.QueryRow(
		"SELECT id, session_id, user_id, username, created_at, expires_at FROM sessions WHERE session_id = ? AND expires_at > ?",
		sessionID, time.Now(),
	).Scan(&session.ID, &session.SessionID, &session.UserID, &session.Username, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found or expired")
		}
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	return &session, nil
}

// DeleteSession 删除会话
func DeleteSession(sessionID string) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// CleanExpiredSessions 清理过期会话
func CleanExpiredSessions() error {
	_, err := DB.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	if err != nil {
		return fmt.Errorf("failed to clean expired sessions: %w", err)
	}
	return nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(userID int) (*User, error) {
	var user User
	err := DB.QueryRow(
		"SELECT id, username, is_admin, created_at, updated_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// GetAllUsers 获取所有用户
func GetAllUsers() ([]User, error) {
	rows, err := DB.Query("SELECT id, username, is_admin, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// CreateUser 创建用户
func CreateUser(username, password string, isAdmin bool) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := DB.Exec(
		"INSERT INTO users (username, password_hash, is_admin) VALUES (?, ?, ?)",
		username, hashedPassword, isAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	return &User{
		ID:       int(id),
		Username: username,
		IsAdmin:  isAdmin,
	}, nil
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(userID int, newPassword string) error {
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = DB.Exec(
		"UPDATE users SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		hashedPassword, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeleteUser 删除用户
func DeleteUser(userID int) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// UpdateUser 更新用户信息
func UpdateUser(userID int, username string, isAdmin bool) error {
	_, err := DB.Exec(
		"UPDATE users SET username = ?, is_admin = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		username, isAdmin, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

