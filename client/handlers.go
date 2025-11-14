package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// LoginPage 登录页面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login - SimpleDBWeb",
	})
}

// LoginAPI 登录API
func LoginAPI(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
			"message": "Username and password are required",
		})
		return
	}

	user, err := AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "invalid_credentials",
			"message": err.Error(),
		})
		return
	}

	session, err := CreateSession(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to create session",
		})
		return
	}

	// 设置Cookie
	c.SetCookie("session_id", session.SessionID, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": gin.H{
			"username": user.Username,
			"is_admin": user.IsAdmin,
		},
	})
}

// LogoutAPI 退出登录API
func LogoutAPI(c *gin.Context) {
	sessionID, _ := c.Get("session_id")
	if sessionIDStr, ok := sessionID.(string); ok {
		DeleteSession(sessionIDStr)
	}

	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

// GetCurrentUserAPI 获取当前用户信息
func GetCurrentUserAPI(c *gin.Context) {
	userID, _ := c.Get("user_id")
	// username, _ := c.Get("username")

	user, err := GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to get user info",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"is_admin": user.IsAdmin,
		},
	})
}

// GetUsersAPI 获取所有用户（管理员）
func GetUsersAPI(c *gin.Context) {
	users, err := GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to get users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// CreateUserAPI 创建用户（管理员）
func CreateUserAPI(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		IsAdmin  bool   `json:"is_admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
			"message": "Username and password are required",
		})
		return
	}

	user, err := CreateUser(req.Username, req.Password, req.IsAdmin)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "user_exists",
				"message": "Username already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User created successfully",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"is_admin": user.IsAdmin,
		},
	})
}

// UpdateUserAPI 更新用户（管理员）
func UpdateUserAPI(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
			"message": "Invalid user ID",
		})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		IsAdmin  bool   `json:"is_admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
			"message": "Username is required",
		})
		return
	}

	if err := UpdateUser(userID, req.Username, req.IsAdmin); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "user_exists",
				"message": "Username already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
	})
}

// DeleteUserAPI 删除用户（管理员）
func DeleteUserAPI(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
			"message": "Invalid user ID",
		})
		return
	}

	// 不能删除自己
	currentUserID, _ := c.Get("user_id")
	if userID == currentUserID.(int) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
			"message": "Cannot delete yourself",
		})
		return
	}

	if err := DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// UpdatePasswordAPI 修改密码
func UpdatePasswordAPI(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
			"message": "Old password and new password are required",
		})
		return
	}

	// 验证旧密码
	_, err := GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to get user info",
		})
		return
	}

	var passwordHash string
	err = DB.QueryRow("SELECT password_hash FROM users WHERE id = ?", userID).Scan(&passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to verify password",
		})
		return
	}

	if !checkPassword(passwordHash, req.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "invalid_password",
			"message": "Old password is incorrect",
		})
		return
	}

	// 更新密码
	if err := UpdateUserPassword(userID.(int), req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal_error",
			"message": "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password updated successfully",
	})
}

