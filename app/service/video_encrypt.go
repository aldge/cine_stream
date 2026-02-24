package service

import (
	"context"
	"errors"
	"time"

	"gitlab.com/cinemae/cine_stream/app/dao"
	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/logger"
)

// VideoEncrypt 视频加密信息业务逻辑
type VideoEncrypt struct {
	ctx             context.Context
	daoVideoEncrypt *dao.VideoEncrypt
}

// NewVideoEncrypt 创建视频加密信息业务逻辑对象
func NewVideoEncrypt(ctx context.Context) *VideoEncrypt {
	return &VideoEncrypt{
		ctx:             ctx,
		daoVideoEncrypt: dao.NewVideoEncrypt(ctx),
	}
}

// Create 保存视频加密信息（单个）
func (v *VideoEncrypt) Create(videoID, key, iv string) error {
	if videoID == "" {
		return errors.New("视频ID不能为空")
	}
	if key == "" {
		return errors.New("加密密钥不能为空")
	}

	encrypt := &entity.VideoEncryptEntity{
		VideoID:    videoID,
		Key:        key,
		IV:         iv,
		CreateTime: uint64(time.Now().Unix()),
	}

	err := v.daoVideoEncrypt.Insert(encrypt)
	if err != nil {
		logger.WithContext(v.ctx).Errorf("[VideoEncrypt.SaveEncryptInfo] 保存视频加密信息失败: %v", err)
		return errors.New("保存视频加密信息失败")
	}

	logger.WithContext(v.ctx).Infof("[VideoEncrypt.SaveEncryptInfo] 保存视频加密信息成功, video_id: %s", videoID)
	return nil
}

// GetEncryptInfoByVideoID 根据video_id获取视频加密信息（兼容旧接口，返回第一个）
func (v *VideoEncrypt) GetEncryptInfoByVideoID(videoID string) (*entity.VideoEncryptEntity, error) {
	if videoID == "" {
		return nil, errors.New("视频ID不能为空")
	}

	encryptInfo, err := v.daoVideoEncrypt.GetByVideoID(videoID)
	if err != nil {
		logger.WithContext(v.ctx).Errorf("[VideoEncrypt.GetEncryptInfoByVideoID] 查询视频加密信息失败: %v", err)
		return nil, errors.New("查询视频加密信息失败")
	}
	return encryptInfo, nil
}

// DeleteEncryptInfoByVideoID 删除指定视频的加密信息
func (v *VideoEncrypt) DeleteEncryptInfoByVideoID(videoID string) error {
	if videoID == "" {
		return errors.New("视频ID不能为空")
	}

	err := v.daoVideoEncrypt.DeleteByVideoID(videoID)
	if err != nil {
		logger.WithContext(v.ctx).Errorf("[VideoEncrypt.DeleteEncryptInfoByVideoID] 删除视频加密信息失败: %v", err)
		return errors.New("删除视频加密信息失败")
	}

	logger.WithContext(v.ctx).Infof("[VideoEncrypt.DeleteEncryptInfoByVideoID] 删除视频加密信息成功, video_id: %s", videoID)
	return nil
}
