package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/app/service"
	"gitlab.com/cinemae/cine_stream/logger"
)

// VideoTsSave 保存视频的 ts 切片
func VideoTsSave(ctx *gin.Context) error {
	var req entity.VideoTSSaveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.WithContext(ctx).Warnf("[VideoTsSave] 参数绑定失败: %v", err)
		return RespJsonError(ctx, 1001, "参数绑定失败")
	}
	if req.VideoID == "" {
		logger.WithContext(ctx).Warnf("[VideoTsSave] 视频ID不能为空")
		return RespJsonError(ctx, 1001, "视频ID不能为空")
	}
	if req.Key == "" {
		logger.WithContext(ctx).Warnf("[VideoTsSave] 视频加密Key不能为空")
		return RespJsonError(ctx, 1001, "视频加密Key不能为空")
	}
	if req.IV == "" {
		logger.WithContext(ctx).Warnf("[VideoTsSave] 视频加密向量不能为空")
		return RespJsonError(ctx, 1001, "视频加密向量不能为空")
	}
	if len(req.TSData) == 0 {
		logger.WithContext(ctx).Warnf("[VideoTsSave] TS列表不能为空")
		return RespJsonError(ctx, 1001, "TS列表不能为空")
	}

	// 验证每个切片的加密信息
	for _, ts := range req.TSData {
		if ts.TSPath == "" {
			logger.WithContext(ctx).Warnf("[VideoTsSave] TS切片path不能为空, ts: %+v", ts)
			return RespJsonError(ctx, 1001, fmt.Sprintf("TS切片path不能为空, ts: %+v", ts))
		}
		if ts.TSSequence < 0 {
			logger.WithContext(ctx).Warnf("[VideoTsSave] TS切片序号不能为空, ts: %+v", ts)
			return RespJsonError(ctx, 1001, fmt.Sprintf("TS切片序号不能为空, ts: %+v", ts))
		}
	}

	// 保存TS切片
	tsService := service.NewVideoTS(ctx)
	err := tsService.BatchCreate(req.VideoID, req.TSData)
	if err != nil {
		logger.WithContext(ctx).Errorf("[VideoTsSave] 批量保存TS切片失败: %v", err)
		return RespJsonError(ctx, 1002, "批量保存TS切片失败")
	}

	// 保存每个切片的加密信息
	encryptService := service.NewVideoEncrypt(ctx)
	err = encryptService.Create(req.VideoID, req.Key, req.IV)
	if err != nil {
		logger.WithContext(ctx).Errorf("[VideoTsSave] 保存视频加密信息失败: %v", err)
		return RespJsonError(ctx, 1003, "保存视频加密信息失败")
	}

	logger.WithContext(ctx).Infof("[VideoTsSave] 批量保存TS切片成功, video_id: %s, count: %d", req.VideoID, len(req.TSData))
	return RespJsonSuccess(ctx, map[string]interface{}{
		"video_id": req.VideoID,
	})
}

// VideoTsList 获取TS切片列表
func VideoTsList(ctx *gin.Context) error {
	videoID := GetParamString(ctx, "video_id")
	definitions := GetParamString(ctx, "definitions")

	if videoID == "" {
		logger.WithContext(ctx).Warnf("[VideoTSList] 视频ID不能为空")
		return RespJsonError(ctx, 1001, "视频ID不能为空")
	}

	// 检查播放权限
	if !service.CheckPlayRights(ctx, videoID) {
		logger.WithContext(ctx).Warnf("[VideoTsList] 用户无播放权限, video_id: %s", videoID)
		ctx.JSON(http.StatusForbidden, &entity.Response{
			Code:    403,
			Message: "无播放权限",
			Data:    make(map[string]interface{}),
		})
		return nil
	}

	tsService := service.NewVideoTS(ctx)

	tsList, err := tsService.GetList(videoID, definitions)
	if err != nil {
		logger.WithContext(ctx).Errorf("[VideoTSList] 查询TS切片列表失败: %v", err)
		return RespJsonError(ctx, 1002, "查询TS切片列表失败")
	}

	return RespJsonSuccess(ctx, map[string]interface{}{
		"ts_list":     tsList,
		"video_id":    videoID,
		"definitions": definitions,
	})
}
