package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许登录页面和登录API访问
		path := c.Request.URL.Path
		if path == "/login" || path == "/api/auth/login" || strings.HasPrefix(path, "/static/") {
			c.Next()
			return
		}

		// 从Cookie获取session ID
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			// 如果请求的是API，返回JSON错误
			if strings.HasPrefix(path, "/api/") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "unauthorized",
					"message": "Please login first",
				})
				c.Abort()
				return
			}
			// 否则重定向到登录页
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 验证session
		session, err := GetSession(sessionID)
		if err != nil {
			// 清除无效的cookie
			c.SetCookie("session_id", "", -1, "/", "", false, true)
			if strings.HasPrefix(path, "/api/") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "unauthorized",
					"message": "Session expired, please login again",
				})
				c.Abort()
				return
			}
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", session.UserID)
		c.Set("username", session.Username)
		c.Set("session_id", sessionID)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "Please login first",
			})
			c.Abort()
			return
		}

		user, err := GetUserByID(userID.(int))
		if err != nil || !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "forbidden",
				"message": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

