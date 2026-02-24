package entity

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/cinemae/cine_stream/consts"
	"gitlab.com/cinemae/cine_stream/utils"
)

// ContextWithRequestID context 添加请求 ID
func ContextWithRequestID(ctx context.Context, loginToken string) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(consts.BizContextKeyRequestID, loginToken)
		return ctx
	}
	return context.WithValue(ctx, consts.BizContextKeyRequestID, loginToken)
}

// ContextValueRequestID context 添加请求 ID
func ContextValueRequestID(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetString(consts.BizContextKeyRequestID)
	}
	loginToken := ctx.Value(consts.BizContextKeyRequestID)
	return fmt.Sprintf("%v", loginToken)
}

// ContextWithLoginAccountID context 添加登录账号ID
func ContextWithLoginAccountID(ctx context.Context, accountID int64) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(consts.BizContextKeyLoginAccountID, accountID)
		return ctx
	}
	return context.WithValue(ctx, consts.BizContextKeyLoginAccountID, accountID)
}

// ContextValueLoginAccountID context 获取登录账号ID
func ContextValueLoginAccountID(ctx context.Context) int64 {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetInt64(consts.BizContextKeyLoginAccountID)
	}
	accountIDIf := ctx.Value(consts.BizContextKeyLoginAccountID)
	accountID, _ := utils.Convert.ToInt64(accountIDIf)
	return accountID
}

// ContextWithLoginAccountName context 添加登录账号名
func ContextWithLoginAccountName(ctx context.Context, accountName string) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(consts.BizContextKeyLoginAccountName, accountName)
		return ctx
	}
	return context.WithValue(ctx, consts.BizContextKeyLoginAccountName, accountName)
}

// ContextValueLoginAccountID context 获取登录账号ID
func ContextValueLoginAccountName(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetString(consts.BizContextKeyLoginAccountName)
	}
	accountName := ctx.Value(consts.BizContextKeyLoginAccountName)
	return fmt.Sprintf("%v", accountName)
}

// ContextWithLoginToken context 添加登录 token
func ContextWithLoginToken(ctx context.Context, loginToken string) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(consts.BizContextKeyLoginToken, loginToken)
		return ctx
	}
	return context.WithValue(ctx, consts.BizContextKeyLoginToken, loginToken)
}

// ContextValueLoginToken context 获取登录 token
func ContextValueLoginToken(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetString(consts.BizContextKeyLoginToken)
	}
	loginToken := ctx.Value(consts.BizContextKeyLoginToken)
	return fmt.Sprintf("%v", loginToken)
}

// ContextWithApplicationName context 添加应用名称
func ContextWithApplicationName(ctx context.Context, applicationName string) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(consts.BizContextKeyApplicationName, applicationName)
		return ctx
	}
	return context.WithValue(ctx, consts.BizContextKeyApplicationName, applicationName)
}

// ContextValueApplicationName context 获取应用名称
func ContextValueApplicationName(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.GetString(consts.BizContextKeyApplicationName)
	}
	applicationName := ctx.Value(consts.BizContextKeyApplicationName)
	return fmt.Sprintf("%v", applicationName)
}
