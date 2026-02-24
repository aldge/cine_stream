package service

import (
	"context"
	"errors"
	"time"

	"gitlab.com/cinemae/cine_stream/app/dao"
	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/logger"
)

// VideoTS TS切片业务逻辑
type VideoTS struct {
	ctx        context.Context
	daoVideoTS *dao.VideoTS
}

// NewVideoTS 创建TS切片业务逻辑对象
func NewVideoTS(ctx context.Context) *VideoTS {
	return &VideoTS{
		ctx:        ctx,
		daoVideoTS: dao.NewVideoTS(ctx),
	}
}

// BatchCreate 批量创建TS切片
func (v *VideoTS) BatchCreate(videoID string, tsList []*entity.VideoTsSaveDataItem) error {
	if videoID == "" {
		return errors.New("视频ID不能为空")
	}
	if len(tsList) == 0 {
		return errors.New("TS切片列表不能为空")
	}
	timeNow := time.Now().Unix()

	var tsEntityList []*entity.VideoTSEntity

	// 设置每个TS切片的默认值
	for _, ts := range tsList {
		if ts.TSSequence < 0 {
			return errors.New("TS序号不能为负数")
		}
		if ts.Duration <= 0 {
			return errors.New("TS时长必须大于0")
		}
		var tsEntity entity.VideoTSEntity
		tsEntity.VideoID = videoID
		tsEntity.TSPath = ts.TSPath
		tsEntity.TSSequence = ts.TSSequence
		tsEntity.Duration = ts.Duration
		tsEntity.Definition = ts.Definition
		tsEntity.CreateTime = timeNow
		tsEntityList = append(tsEntityList, &tsEntity)
	}

	// 批量保存到数据库
	err := v.daoVideoTS.BatchInsert(videoID, tsEntityList)
	if err != nil {
		logger.WithContext(v.ctx).Errorf("[VideoTS.BatchCreate] 批量保存TS切片失败: %v", err)
		return err
	}

	logger.WithContext(v.ctx).Infof("[VideoTS.BatchCreate] 批量保存TS切片成功, video_id: %s, count: %d", videoID, len(tsList))
	return nil
}

// GetList 获取视频的TS切片列表
func (v *VideoTS) GetList(videoID string, definitions string) ([]entity.VideoTSEntity, error) {
	if videoID == "" {
		return nil, errors.New("视频ID不能为空")
	}
	tsList, err := v.daoVideoTS.GetByVideoID(videoID, definitions)
	if err != nil {
		logger.WithContext(v.ctx).Errorf("[VideoTS.GetList] 查询TS切片列表失败: %v", err)
		return nil, errors.New("查询TS切片列表失败")
	}
	return tsList, nil
}
