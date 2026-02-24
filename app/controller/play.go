package controller

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/app/service"
	"gitlab.com/cinemae/cine_stream/logger"
	"gitlab.com/cinemae/cine_stream/utils"
	"gitlab.com/cinemae/gopkg/app"
)

const (
	cinePlayerKey    = "0123456789abcdef" // 16字节
	cinePlayerNonce  = "0123456789ab"     // 12字节
	cinePlayerNonce2 = "0123456789ac"     // 12字节
)

// Play 获取播放的 m3u8 文件
func Play(ctx *gin.Context) error {
	videoID := ctx.Param("video_id")
	if videoID == "" {
		logger.WithContext(ctx).Warnf("[PlayIndexM3u8] 视频ID不能为空")
		return RespJsonError(ctx, 1001, "视频ID不能为空")
	}

	// 检查播放权限
	if !service.CheckPlayRights(ctx, videoID) {
		logger.WithContext(ctx).Warnf("[Play] 用户无播放权限, video_id: %s", videoID)
		ctx.JSON(http.StatusForbidden, &entity.Response{
			Code:    403,
			Message: "无播放权限",
			Data:    make(map[string]interface{}),
		})
		return nil
	}

	appName := app.GetAppName(ctx)

	// 生成主 m3u8
	m3u8Content := `#EXTM3U
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=4096000,RESOLUTION=1920x1080
/play/%s/index.m3u8?app=%s`
	ctx.Header("Content-Type", "application/vnd.apple.mpegurl")
	ctx.String(http.StatusOK, fmt.Sprintf(m3u8Content, videoID, appName))

	return nil
}

// PlayHlsIndexM3u8 获取播放的 hls m3u8 文件
func PlayHlsIndexM3u8(ctx *gin.Context) error {

	videoID := ctx.Param("video_id")
	if videoID == "" {
		logger.WithContext(ctx).Warnf("[PlayHlsIndexM3u8] 视频ID不能为空")
		return RespJsonError(ctx, 1001, "视频ID不能为空")
	}

	// 检查播放权限
	if !service.CheckPlayRights(ctx, videoID) {
		logger.WithContext(ctx).Warnf("[PlayHlsIndexM3u8] 用户无播放权限, video_id: %s", videoID)
		ctx.JSON(http.StatusForbidden, &entity.Response{
			Code:    403,
			Message: "无播放权限",
			Data:    make(map[string]interface{}),
		})
		return nil
	}

	// 获取 app 参数，确保中间件验证通过（虽然 service 层也会获取，但这里显式获取以确保验证）
	_ = app.GetAppName(ctx)

	// 获取所有的 ts 分片
	tsService := service.NewVideoTS(ctx)
	tsList, err := tsService.GetList(videoID, "")
	if err != nil {
		logger.WithContext(ctx).Errorf("[PlayHlsIndexM3u8] 查询TS切片列表失败: %v", err)
		return RespJsonError(ctx, 1002, "查询TS切片列表失败")
	}

	playService := service.NewPlay(ctx)
	m3u8Content, err := playService.GenerateM3U8Content(ctx, videoID, tsList)
	if err != nil {
		logger.WithContext(ctx).Errorf("[PlayHlsIndexM3u8] 生成m3u8内容失败: %v", err)
		return RespJsonError(ctx, 1003, "生成m3u8内容失败")
	}

	// 返回给前端 m3u8 文件
	ctx.Header("Content-Type", "application/vnd.apple.mpegurl")
	ctx.String(http.StatusOK, m3u8Content)
	return nil
}

// PlayHlsIndexEncKey 获取播放的 hls 加密 key（通过 video_encrypt_id）
func PlayHlsIndexEncKey(ctx *gin.Context) error {

	videIDStr := ctx.Param("video_id")
	if videIDStr == "" {
		logger.WithContext(ctx).Warnf("[PlayHlsIndexEncKey] video_encrypt_id 不能为空")
		return RespJsonError(ctx, 1001, "video_encrypt_id 不能为空")
	}

	// 检查播放权限
	if !service.CheckPlayRights(ctx, videIDStr) {
		logger.WithContext(ctx).Warnf("[PlayHlsIndexEncKey] 用户无播放权限, video_id: %s", videIDStr)
		ctx.JSON(http.StatusForbidden, &entity.Response{
			Code:    403,
			Message: "无播放权限",
			Data:    make(map[string]interface{}),
		})
		return nil
	}

	// 获取视频加密信息
	encryptService := service.NewVideoEncrypt(ctx)
	encrypt, err := encryptService.GetEncryptInfoByVideoID(videIDStr)
	if err != nil {
		logger.WithContext(ctx).Errorf("[PlayHlsIndexEncKey] 获取视频加密信息失败: %v", err)
		return RespJsonError(ctx, 1002, "获取视频加密信息失败")
	}

	// 返回前端的 application/octet-stream
	ctx.Header("Content-Type", "application/octet-stream")
	// 返回给前端加密 key
	keyBytes, err := hex.DecodeString(encrypt.Key)
	if err != nil {
		logger.WithContext(ctx).Errorf("[PlayHlsIndexEncKey] 解码加密 key 失败: %v", err)
		return RespJsonError(ctx, 1003, "解码加密 key 失败")
	}
	ctx.Data(http.StatusOK, "application/octet-stream", keyBytes)
	return nil
}

// PlayCine cine 播放器协议播放接口（私有协议暂时不用一级m3u8）
// func PlayCine(ctx *gin.Context) error {
// 	videoID := ctx.Param("video_id")
// 	if videoID == "" {
// 		logger.WithContext(ctx).Warnf("[PlayIndexM3u8] 视频ID不能为空")
// 		return RespJsonError(ctx, 1001, "视频ID不能为空")
// 	}

// 	// 生成主 m3u8
// 	m3u8Content := `#EXTM3U
// #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=4096000,RESOLUTION=1920x1080
// %s/play/cine/%s/index.m3u8`

// 	baseUrl := utils.GetRequestBaseURL(ctx)
// 	m3u8Content = fmt.Sprintf(m3u8Content, baseUrl, videoID)

// 	// m3u8 数据加密为二进制数据， GCM模式示例
// 	m3u8Data, _, err := utils.AESEncrypt(m3u8Content, []byte(cinePlayerKey), []byte(cinePlayerNonce), utils.ModeGCM)
// 	if err != nil {
// 		logger.WithContext(ctx).Errorf("[PlayCine] GCM加密失败：%v", err)
// 		return RespJsonError(ctx, 1003, "获取视频信息失败 crypto err")
// 	}
// 	m3u8base64Str := base64.StdEncoding.EncodeToString(m3u8Data)
// 	var result = map[string]interface{}{
// 		"info": m3u8base64Str,
// 	}
// 	RespJsonSuccess(ctx, result)
// 	return nil
// }

// PlayCineHlsIndexC3u8 获取cine播放器的 hls m3u8 文件
func PlayCineHlsIndexC3u8(ctx *gin.Context) error {

	videoID := ctx.Param("video_id")
	if videoID == "" {
		logger.WithContext(ctx).Warnf("[PlayCineHlsIndexM3u8] 视频ID不能为空")
		return RespJsonError(ctx, 1001, "视频ID不能为空")
	}

	// 检查播放权限
	if !service.CheckPlayRights(ctx, videoID) {
		logger.WithContext(ctx).Warnf("[PlayCineHlsIndexC3u8] 用户无播放权限, video_id: %s", videoID)
		ctx.JSON(http.StatusForbidden, &entity.Response{
			Code:    403,
			Message: "无播放权限",
			Data:    make(map[string]interface{}),
		})
		return nil
	}

	// 获取所有的 ts 分片
	tsService := service.NewVideoTS(ctx)
	tsList, err := tsService.GetList(videoID, "")
	if err != nil {
		logger.WithContext(ctx).Errorf("[PlayCineHlsIndexM3u8] 查询TS切片列表失败: %v", err)
		return RespJsonError(ctx, 1002, "查询TS切片列表失败")
	}

	playService := service.NewPlay(ctx)
	m3u8Content, err := playService.GenerateM3U8Content(ctx, videoID, tsList)
	if err != nil {
		logger.WithContext(ctx).Errorf("[PlayCineHlsIndexM3u8] 生成m3u8内容失败: %v", err)
		return RespJsonError(ctx, 1003, "生成m3u8内容失败")
	}

	// m3u8 数据加密为二进制数据， GCM模式示例
	m3u8Data, _, err := utils.AESEncrypt(m3u8Content, []byte(cinePlayerKey), []byte(cinePlayerNonce), utils.ModeGCM)
	if err != nil {
		logger.WithContext(ctx).Errorf("[PlayCineHlsIndexM3u8] GCM加密失败：%v", err)
		return RespJsonError(ctx, 1003, "获取视频信息失败 crypto err")
	}
	m3u8base64Str := base64.StdEncoding.EncodeToString(m3u8Data)
	var result = map[string]interface{}{
		"info": m3u8base64Str,
	}

	return RespJsonSuccess(ctx, result)
}
