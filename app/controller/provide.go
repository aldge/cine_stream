package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/app/service"
	"gitlab.com/cinemae/cine_stream/logger"
)

// ProvideIndex 资源提供接口
func ProvideIndex(ctx *gin.Context) error {
	action := GetParamString(ctx, "ac")
	tag := GetParamString(ctx, "t")
	hour := GetParamString(ctx, "h")
	ids := GetParamString(ctx, "ids")
	word := GetParamString(ctx, "wd")

	// 解析分页参数
	page := GetParamIntDef(ctx, "pg", 1)
	limit := GetParamIntDef(ctx, "limit", 20)
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	provideService := service.NewProvideService(ctx)

	var result interface{}
	var err error
	var logMsg string

	// 如果没有传入任何参数（action为空），返回简化列表
	if action == "" || action == "list" {
		result, err = provideService.GetSimpleVideoList(page, limit, tag, hour, ids, word)
		logMsg = "[ProvideIndex] 获取简化视频列表失败"
	} else {
		// 如果有 ac 参数，返回完整列表
		result, err = provideService.GetFullVideoList(page, limit, tag, hour, ids, word)
		logMsg = "[ProvideIndex] 获取完整视频列表失败"
	}

	if err != nil {
		logger.WithContext(ctx).Errorf("%s: %v", logMsg, err)
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "获取数据失败",
			"list": []interface{}{},
		})
		return nil
	}

	ctx.JSON(http.StatusOK, result)
	return nil
}

// ProvideSave 保存视频信息接口
func ProvideSave(ctx *gin.Context) error {
	var vodList []entity.VodEntity
	if err := ctx.ShouldBindJSON(&vodList); err != nil {
		logger.WithContext(ctx).Warnf("[ProvideSave] 参数绑定失败: %v", err)
		return RespJsonError(ctx, 1001, "参数绑定失败")
	}

	if len(vodList) == 0 {
		logger.WithContext(ctx).Warnf("[ProvideSave] 视频列表为空")
		return RespJsonError(ctx, 1001, "视频列表不能为空")
	}

	// 转换为指针数组
	vodPtrList := make([]*entity.VodEntity, 0, len(vodList))
	for i := range vodList {
		vodPtrList = append(vodPtrList, &vodList[i])
	}

	provideService := service.NewProvideService(ctx)
	err := provideService.BatchSave(vodPtrList)
	if err != nil {
		logger.WithContext(ctx).Errorf("[ProvideSave] 批量保存视频信息失败: %v", err)
		return RespJsonError(ctx, 1002, "批量保存视频信息失败")
	}

	// 收集所有保存成功的 vod_id
	vodIDs := make([]int64, 0, len(vodList))
	for _, vod := range vodList {
		if vod.VodID > 0 {
			vodIDs = append(vodIDs, vod.VodID)
		}
	}

	logger.WithContext(ctx).Infof("[ProvideSave] 批量保存视频信息成功, count: %d, vod_ids: %v", len(vodList), vodIDs)
	return RespJsonSuccess(ctx, map[string]interface{}{
		"count":   len(vodList),
		"vod_ids": vodIDs,
	})
}
