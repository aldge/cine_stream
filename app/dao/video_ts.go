package dao

import (
	"context"

	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/config"
	"gorm.io/gorm"
)

const (
	videoTsDBName          = "cine_stream"   // TS 切片表数据库名，不同的 app 以 dbName_appName 形式命名
	videoTsTabelNamePrefix = "cine_video_ts" // TS切片表名名称
)

// VideoTS TS切片数据访问对象
type VideoTS struct {
	ctx context.Context
	db  *gorm.DB
}

// NewVideoTS 创建TS切片数据访问对象
func NewVideoTS(ctx context.Context) *VideoTS {
	videoTs := &VideoTS{
		ctx: ctx,
	}
	dbName := getAppDBName(ctx, videoTsDBName)
	videoTs.db = GetDB(dbName)
	// 如果找不到带 app 后缀的数据库配置，回退到默认数据库配置
	if videoTs.db == nil && dbName != videoTsDBName {
		videoTs.db = GetDB(videoTsDBName)
	}
	return videoTs
}

// getTableName 获取表名
func (vs *VideoTS) getTableName(videoID string) string {
	dbName := getAppDBName(vs.ctx, videoTsDBName)
	tableConf := config.GetAppConf().GetDatabaseTableConf(dbName, videoTsTabelNamePrefix)
	// 如果找不到带 app 后缀的数据库配置，使用默认数据库配置
	if tableConf.ShardingNum == 0 && dbName != videoTsDBName {
		tableConf = config.GetAppConf().GetDatabaseTableConf(videoTsDBName, videoTsTabelNamePrefix)
	}
	return getShardingTableName(videoID, tableConf.ShardingNum, videoTsTabelNamePrefix)
}

// BatchInsert 批量插入TS切片记录
func (vs *VideoTS) BatchInsert(videoID string, tsList []*entity.VideoTSEntity) error {
	if len(tsList) == 0 {
		return ErrInvalidParam
	}
	if vs.db == nil {
		return ErrDBConfNotFound
	}
	return vs.db.Table(vs.getTableName(videoID)).CreateInBatches(tsList, 100).Error
}

// GetByVideoID 根据视频ID获取TS切片列表
func (vs *VideoTS) GetByVideoID(videoID string, definitions string) ([]entity.VideoTSEntity, error) {
	if videoID == "" {
		return nil, ErrInvalidParam
	}
	if vs.db == nil {
		return nil, ErrDBConfNotFound
	}

	var tsList []entity.VideoTSEntity
	db := vs.db.Table(vs.getTableName(videoID)).Where("video_id = ?", videoID)

	// 按清晰度过滤
	if definitions != "" {
		db = db.Where("definition = ?", definitions)
	}

	err := db.Order("ts_sequence ASC").Find(&tsList).Error
	if err != nil {
		return nil, err
	}
	return tsList, nil
}

// DeleteByVideoID 删除指定视频的所有TS切片
func (vs *VideoTS) DeleteByVideoID(videoID string) error {
	if videoID == "" {
		return ErrInvalidParam
	}
	if vs.db == nil {
		return ErrDBConfNotFound
	}
	return vs.db.Table(vs.getTableName(videoID)).Where("video_id = ?", videoID).Delete(&entity.VideoTSEntity{}).Error
}

// GetCountByVideoID 获取指定视频的TS切片数量
func (vs *VideoTS) GetCountByVideoID(videoID string) (int64, error) {
	if videoID == "" {
		return 0, ErrInvalidParam
	}
	if vs.db == nil {
		return 0, ErrDBConfNotFound
	}

	var count int64
	err := vs.db.Table(vs.getTableName(videoID)).Where("video_id = ?", videoID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
