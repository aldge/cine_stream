package controller

import (
	"fmt"
	"gitlab.com/cinemae/cine_stream/utils"
	"github.com/gin-gonic/gin"
)
type ginCtx gin.Context

// 一些基础的公共方法

// GetParamString 获取参数返回 string
func GetParamString(ctx *gin.Context, key string) string {
	queryVal, ok := ctx.GetQuery(key)
	if ok {
		return queryVal
	}
	return ctx.PostForm(key)
}

// GetParamStringDef 获取参数返回 string 带默认值
func GetParamStringDef(ctx *gin.Context, key string, def string) string {
	queryVal, ok := ctx.GetQuery(key)
	if ok {
		return queryVal
	}
	return ctx.DefaultPostForm(key, def)
}

// GetParamInt 获取参数返回 int
func GetParamInt(ctx *gin.Context, key string) int {
	queryVal := GetParamString(ctx, key)
	return utils.Convert.StringToInt(queryVal)
}

// GetParamIntDef 获取参数返回 int
func GetParamIntDef(ctx *gin.Context, key string, def int) int {
	queryVal := GetParamStringDef(ctx, key, fmt.Sprintf("%d", def))
	return utils.Convert.StringToInt(queryVal)
}

// GetParamInt64 获取参数返回 float64
func GetParamInt64(ctx *gin.Context, key string) int64 {
	queryVal := GetParamString(ctx, key)
	return utils.Convert.StringToInt64(queryVal)
}

// GetParamInt64Def 获取参数返回 float64
func GetParamInt64Def(ctx *gin.Context, key string, def int64) int64 {
	queryVal := GetParamStringDef(ctx, key, fmt.Sprintf("%d", def))
	return utils.Convert.StringToInt64(queryVal)
}

// GetParamFloat32 获取参数返回 float32
func GetParamFloat32(ctx *gin.Context, key string) float32 {
	queryVal := GetParamString(ctx, key)
	return utils.Convert.StringToFloat32(queryVal)
}

// GetParamFloat32Def 获取参数返回 float32
func GetParamFloat32Def(ctx *gin.Context, key string, def float32) float32 {
	queryVal := GetParamStringDef(ctx, key, fmt.Sprintf("%f", def))
	return utils.Convert.StringToFloat32(queryVal)
}

// GetParamFloat64 获取参数返回 float64
func GetParamFloat64(ctx *gin.Context, key string) float64 {
	queryVal := GetParamString(ctx, key)
	return utils.Convert.StringToFloat64(queryVal)
}

// GetParamFloat64Def 获取参数返回 float64
func GetParamFloat64Def(ctx *gin.Context, key string, def float64) float64 {
	queryVal := GetParamStringDef(ctx, key, fmt.Sprintf("%f", def))
	return utils.Convert.StringToFloat64(queryVal)
}

// GetParamMap 获取参数返回字符串数组
func GetParamMap(ctx *gin.Context, key string) map[string]string {
	queryVal, ok := ctx.GetQueryMap(key)
	if ok {
		return queryVal
	}
	return ctx.PostFormMap(key)
}