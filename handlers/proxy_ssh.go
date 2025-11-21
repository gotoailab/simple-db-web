package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/gotoailab/simple-db-web/database"
	"golang.org/x/crypto/ssh"
)

// SSHProxyConfig SSH代理配置
type SSHProxyConfig struct {
	Host     string `json:"host"`     // SSH服务器地址
	Port     string `json:"port"`     // SSH服务器端口，默认22
	User     string `json:"user"`     // SSH用户名
	Password string `json:"password"` // SSH密码（如果使用密码认证）
	KeyFile  string `json:"key_file"` // SSH私钥文件路径（如果使用密钥认证）
	KeyData  string `json:"key_data"` // SSH私钥内容（base64编码，如果使用密钥认证）
}

// SSHProxy SSH代理实现
type SSHProxy struct {
	client *ssh.Client
	config *SSHProxyConfig
}

// NewSSHProxy 创建SSH代理
// config: 代理配置的JSON字符串（database.ProxyConfig 的 JSON）
func NewSSHProxy(config string) (Proxy, error) {
	// 先解析为 database.ProxyConfig，因为前端发送的是这个结构
	var dbProxyConfig database.ProxyConfig
	if err := json.Unmarshal([]byte(config), &dbProxyConfig); err != nil {
		return nil, fmt.Errorf("解析SSH代理配置失败: %w", err)
	}

	// 转换为 SSHProxyConfig
	// 注意：密码和私钥已经在 Connect 函数中解密，这里直接使用
	proxyConfig := SSHProxyConfig{
		Host:     dbProxyConfig.Host,
		Port:     dbProxyConfig.Port,
		User:     dbProxyConfig.User,
		Password: dbProxyConfig.Password, // 已经是解密后的密码
		KeyFile:  dbProxyConfig.KeyFile,
	}

	// 从 Config 字段中提取 key_data（如果存在）
	if dbProxyConfig.Config != "" {
		var configMap map[string]interface{}
		if err := json.Unmarshal([]byte(dbProxyConfig.Config), &configMap); err == nil {
			if keyData, ok := configMap["key_data"].(string); ok && keyData != "" {
				proxyConfig.KeyData = keyData // 已经是解密后的私钥
			}
		}
	}

	// 设置默认端口
	if proxyConfig.Port == "" {
		proxyConfig.Port = "22"
	}

	// 构建SSH客户端配置
	sshConfig := &ssh.ClientConfig{
		User:            proxyConfig.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境应该验证主机密钥
		Timeout:         10 * time.Second,
	}

	// 认证方式：优先使用密钥，其次使用密码
	var authMethods []ssh.AuthMethod

	if proxyConfig.KeyData != "" {
		// 使用提供的密钥数据（已经是解密后的原始私钥内容）
		signer, err := ssh.ParsePrivateKey([]byte(proxyConfig.KeyData))
		if err != nil {
			return nil, fmt.Errorf("解析SSH私钥失败: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if proxyConfig.Password != "" {
		// 使用密码认证（已经是解密后的原始密码）
		authMethods = append(authMethods, ssh.Password(proxyConfig.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("SSH代理需要提供密码或密钥")
	}

	sshConfig.Auth = authMethods

	// 连接到SSH服务器
	address := net.JoinHostPort(proxyConfig.Host, proxyConfig.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("连接SSH服务器失败: %w", err)
	}

	return &SSHProxy{
		client: client,
		config: &proxyConfig,
	}, nil
}

// Dial 通过SSH隧道建立到目标地址的连接
func (s *SSHProxy) Dial(network, address string) (net.Conn, error) {
	if s.client == nil {
		return nil, fmt.Errorf("SSH客户端未初始化")
	}

	// 通过SSH隧道连接到目标地址
	conn, err := s.client.Dial(network, address)
	if err != nil {
		return nil, fmt.Errorf("通过SSH隧道连接失败: %w", err)
	}

	return conn, nil
}

// Close 关闭SSH连接
func (s *SSHProxy) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// buildSSHProxyConfig 从ProxyConfig构建SSH代理配置JSON
func buildSSHProxyConfig(proxyConfig *database.ProxyConfig) (string, error) {
	sshConfig := SSHProxyConfig{
		Host:     proxyConfig.Host,
		Port:     proxyConfig.Port,
		User:     proxyConfig.User,
		Password: proxyConfig.Password, // 已经是解密后的密码
		KeyFile:  proxyConfig.KeyFile,
	}

	// 如果提供了Config字段，尝试解析并合并
	if proxyConfig.Config != "" {
		var extraConfig map[string]interface{}
		if err := json.Unmarshal([]byte(proxyConfig.Config), &extraConfig); err == nil {
			// 合并额外配置
			if keyData, ok := extraConfig["key_data"].(string); ok && keyData != "" {
				// keyData 是加密后的，需要解密（使用 handlers 包中的 decryptPassword 函数）
				decryptedKeyData, err := decryptPassword(keyData)
				if err != nil {
					return "", fmt.Errorf("解密SSH私钥失败: %w", err)
				}
				sshConfig.KeyData = decryptedKeyData
			}
		}
	}

	configJSON, err := json.Marshal(sshConfig)
	if err != nil {
		return "", fmt.Errorf("序列化SSH配置失败: %w", err)
	}

	return string(configJSON), nil
}
