// Package utils 公共工具包
package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetConfPaths 获取配置文件 path
func GetConfPaths(rootPath string) []string {
	return []string{
		rootPath,
		rootPath + "/conf/",
		rootPath + "/../conf/",
		rootPath + "/../../conf/",
	}
}

// GetRequestBaseURL 构建当前请求的绝对域名（含协议和端口）
func GetRequestBaseURL(ctx *gin.Context) string {
	scheme := "http"
	// 优先从 X-Forwarded-Proto 获取代理后的协议
	if forwardedProto := ctx.Request.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	} else if ctx.Request.TLS != nil {
		scheme = "https"
	}

	// 优先从 X-Forwarded-Host 获取代理服务器的 host
	forwardedHost := ctx.Request.Header.Get("X-Forwarded-Host")
	var host string
	// 打印 X-Forwarded-Host
	forwardedPort := ""
	if forwardedHost != "" {
		host = forwardedHost
		// 检查 X-Forwarded-Host 是否已包含端口（包含冒号表示有端口）
		if !strings.Contains(host, ":") {
			// 如果没有端口，尝试从 X-Forwarded-Port 获取端口
			forwardedPort = ctx.Request.Header.Get("X-Forwarded-Port")
			if forwardedPort != "" {
				host = fmt.Sprintf("%s:%s", host, forwardedPort)
			}
		}
	} else {
		// 如果没有代理，使用当前服务器的 host（通常已包含端口）
		host = ctx.Request.Host
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}
