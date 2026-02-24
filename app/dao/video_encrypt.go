package dao

import (
	"context"

	"gitlab.com/cinemae/cine_stream/app/entity"
	"gorm.io/gorm"
)

const (
	videoEncryptTableName = "cine_video_encrypt" // 视频加密信息表名
)

// VideoEncrypt 视频加密信息数据访问对象
type VideoEncrypt struct {
	ctx context.Context
	db  *gorm.DB
}

// NewVideoEncrypt 创建视频加密信息数据访问对象
func NewVideoEncrypt(ctx context.Context) *VideoEncrypt {
	ve := &VideoEncrypt{
		ctx: ctx,
	}
	dbName := getAppDBName(ctx, videoTsDBName)
	ve.db = GetDB(dbName)
	// 如果找不到带 app 后缀的数据库配置，回退到默认数据库配置
	if ve.db == nil && dbName != videoTsDBName {
		ve.db = GetDB(videoTsDBName)
	}
	return ve
}

// Insert 保存视频加密信息（单个）
func (ve *VideoEncrypt) Insert(encrypt *entity.VideoEncryptEntity) error {
	if encrypt.VideoID == "" {
		return ErrInvalidParam
	}
	if ve.db == nil {
		return ErrDBConfNotFound
	}
	return ve.db.Table(videoEncryptTableName).Create(encrypt).Error
}

// BatchInsert 批量保存视频加密信息
func (ve *VideoEncrypt) BatchInsert(encryptList []*entity.VideoEncryptEntity) error {
	if len(encryptList) == 0 {
		return ErrInvalidParam
	}
	if ve.db == nil {
		return ErrDBConfNotFound
	}
	return ve.db.Table(videoEncryptTableName).CreateInBatches(encryptList, 100).Error
}

// GetByVideoID 根据video_id查询视频加密信息列表
func (ve *VideoEncrypt) GetByVideoID(videoID string) (*entity.VideoEncryptEntity, error) {
	if videoID == "" {
		return nil, ErrInvalidParam
	}
	if ve.db == nil {
		return nil, ErrDBConfNotFound
	}
	var encrypt = entity.VideoEncryptEntity{}
	err := ve.db.Table(videoEncryptTableName).Where("video_id = ?", videoID).First(&encrypt).Error
	if err != nil {
		return nil, err
	}
	return &encrypt, nil
}

// DeleteByVideoID 删除指定视频的加密信息
func (ve *VideoEncrypt) DeleteByVideoID(videoID string) error {
	if videoID == "" {
		return ErrInvalidParam
	}
	if ve.db == nil {
		return ErrDBConfNotFound
	}
	return ve.db.Table(videoEncryptTableName).Where("video_id = ?", videoID).Delete(&entity.VideoEncryptEntity{}).Error
}
