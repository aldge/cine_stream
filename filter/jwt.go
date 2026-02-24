package filter

import (
	"reflect"
	"strings"

	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/consts"
	"gitlab.com/cinemae/cine_stream/logger"

	"github.com/gin-gonic/gin"
	passport "gitlab.com/cinemae/gopkg/casdoor"
)

var (
	ignoreLoginAuthPath = map[string]bool{
		"/admin/auth/login": true,
	}
)

// AuthLoginJWT 登录 token 认证
func AuthLoginJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authLoginJWTHandle(c) {
			c.Next()
		} else {
			c.Abort()
		}
	}
}

func authLoginJWTHandle(c *gin.Context) (keepNext bool) {
	var tokenString string

	// 优先从 Authorization header 获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// 解析 Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString = parts[1]
		} else {
			// Authorization header 格式不正确，但不阻止请求，只是不设置用户信息
			logger.WithContext(c).Warnf("[AuthLoginJWT] Invalid Authorization header format")
		}
	}

	// 如果 header 中没有 token，尝试从 cookie 中获取
	if tokenString == "" {
		cookieToken, err := c.Cookie("token")
		if err == nil && cookieToken != "" {
			tokenString = cookieToken
		} else {
			// 只有在确实没有 cookie 时才记录警告
			logger.WithContext(c).Debugf("[CookieJWT] No token cookie found: %v", err)
		}
	}

	// 如果都没有 token，允许继续，但不设置用户信息
	if tokenString == "" {
		return true
	}

	// 解析 JWT token
	claims, err := passport.ParseJwtToken(tokenString)
	if err != nil {
		// Token 解析失败，记录日志但允许继续（可能是不需要登录的接口）
		logger.WithContext(c).Warnf("[AuthLoginJWT] Failed to parse JWT token: %v", err)
		return true
	}

	// 提取用户信息
	userIDStr := claims.Id
	userName := ""

	// 从 claims 中获取用户名，使用反射安全获取 PreferredUsername 字段
	// claims 是指针类型，需要先解引用
	claimsValue := reflect.ValueOf(claims)
	if claimsValue.Kind() == reflect.Ptr {
		claimsValue = claimsValue.Elem()
	}
	if claimsValue.Kind() == reflect.Struct {
		if preferredUsernameField := claimsValue.FieldByName("PreferredUsername"); preferredUsernameField.IsValid() && preferredUsernameField.Kind() == reflect.String {
			userName = preferredUsernameField.String()
		} else if nameField := claimsValue.FieldByName("Name"); nameField.IsValid() && nameField.Kind() == reflect.String {
			userName = nameField.String()
		} else if usernameField := claimsValue.FieldByName("Username"); usernameField.IsValid() && usernameField.Kind() == reflect.String {
			userName = usernameField.String()
		}
	}

	applicationName := claims.Owner // Owner 通常是组织名称，对应应用名称

	// 存储到 context 中
	// 注意：Casdoor 的用户ID是字符串（UUID），直接存储为字符串
	// 使用 gin context 的 Set 方法存储，后续可以通过 GetString 获取
	c.Set(consts.BizContextKeyLoginAccountID, userIDStr)
	entity.ContextWithLoginAccountName(c, userName)
	entity.ContextWithApplicationName(c, applicationName)
	entity.ContextWithLoginToken(c, tokenString)

	logger.WithContext(c).Infof("[AuthLoginJWT] User authenticated: userID=%s, userName=%s, appName=%s", userIDStr, userName, applicationName)

	return true
}
