package filter

import (
	"net/http"
	"time"

	"gitlab.com/cinemae/cine_stream/logger"
	"gitlab.com/cinemae/gopkg/app"
	"gitlab.com/cinemae/gopkg/errors"

	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// DebugCosTime 请求耗时打印
func DebugCosTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		cost := time.Since(start)
		logger.WithContext(c).Infof("[RequestCosTime] path=%+v method=%+v status=%+v cost=%+v",
			path, c.Request.Method, c.Writer.Status(), cost)
	}
}

// RequestParse 请求公共参数解析
func RequestParse() gin.HandlerFunc {
	return func(c *gin.Context) {
		if requestParseHandle(c) {
			c.Next()
		} else {
			c.Abort()
		}
	}
}

func requestParseHandle(c *gin.Context) (keepNext bool) {
	// 404 直接返回，不继续往下执行拦截器
	if c.Writer.Status() == http.StatusNotFound {
		return false
	}

	// Swagger 静态文件路径跳过 app 验证
	path := c.Request.URL.Path
	if path == "/swagger" || path == "/swagger/" || len(path) > 9 && path[:9] == "/swagger/" {
		return true
	}

	// 判断 app 是否合法
	if !app.CtxAppIsValid(c) {
		c.JSON(http.StatusOK, &errResponse{
			Code:    int32(errors.ClientUnknownError),
			Message: "invalid app",
			Data:    make(map[string]interface{}),
		})
		c.Abort()
		return false
	}
	// 检查数据库配置
	return true
}
