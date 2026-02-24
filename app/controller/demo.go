package controller

import "github.com/gin-gonic/gin"

// DemoIndex demo
func DemoIndex(ctx *gin.Context) error {

	return RespJsonSuccess(ctx, "hello, this is cine_stream")
}
