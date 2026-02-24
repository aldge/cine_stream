package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.com/cinemae/cine_stream/config"
	"gitlab.com/cinemae/cine_stream/consts"
	klog "gitlab.com/cinemae/gopkg/log"
)

const (
	appDefaultLog = "default"
)

// Init 初始化日志
func Init() {
	if len(config.GetAppConf().GetLoggerConf()) == 0 {
		return
	}
	for logName, loggerConf := range config.GetAppConf().GetLoggerConf() {
		defaultLog := klog.NewZapLogWithCallerSkip(loggerConf, 2)
		klog.Register(logName, defaultLog)
	}
}

// WithContextFields 带自定义字段到 context
func WithContextFields(ctx context.Context, fields ...string) context.Context {
	if ctx == nil {
		return ctx
	}
	ginCtx, ok := ctx.(*gin.Context)
	if ok && ginCtx != nil {
		ginCtx.Set(consts.BizContextKeyLogger, WithFields(fields...))
		return ctx
	}
	return context.WithValue(ctx, consts.BizContextKeyLogger, WithFields(fields...))
}

// WithContext context logger
func WithContext(ctx context.Context) klog.Logger {
	if ctx == nil {
		return klog.GetLogger(appDefaultLog)
	}
	ginCtx, ok := ctx.(*gin.Context)
	if ok && ginCtx != nil {
		ginCtxLogger, exists := ginCtx.Get(consts.BizContextKeyLogger)
		if exists {
			return ginCtxLogger.(klog.Logger)
		}
	}
	if ctxLogger, ok := ctx.Value(consts.BizContextKeyLogger).(klog.Logger); ok {
		return ctxLogger
	}
	return klog.GetLogger(appDefaultLog)
}

// Fatalf Fatal log
func Fatalf(format string, v ...interface{}) {
	klog.GetLogger(appDefaultLog).Fatalf(format, v...)
}

// CtxErrorf
func Errorf(format string, v ...interface{}) {
	klog.GetLogger(appDefaultLog).Errorf(format, v...)
}

// Warnf Warn log
func Warnf(format string, v ...interface{}) {
	klog.GetLogger(appDefaultLog).Warnf(format, v...)
}

// Infof Info log
func Infof(format string, v ...interface{}) {
	klog.GetLogger(appDefaultLog).Infof(format, v...)
}

// Debugf Debug log
func Debugf(format string, v ...interface{}) {
	klog.GetLogger(appDefaultLog).Debugf(format, v...)
}

// WithFields 额外参数
func WithFields(fields ...string) klog.Logger {
	return klog.GetLogger(appDefaultLog).WithFields(fields...)
}

// Sync 同步
func Sync() {
	klog.GetLogger(appDefaultLog).Sync()
}
