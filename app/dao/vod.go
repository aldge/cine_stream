package dao

import (
	"context"
	"strconv"
	"strings"

	"gitlab.com/cinemae/cine_stream/app/entity"
	"gorm.io/gorm"
)

const (
	vodDBName = "cine_stream" // VOD 表数据库名
)

// Vod VOD数据访问对象
type Vod struct {
	ctx context.Context
	db  *gorm.DB
}

// NewVod 创建VOD数据访问对象
func NewVod(ctx context.Context) *Vod {
	vod := &Vod{
		ctx: ctx,
	}
	dbName := getAppDBName(ctx, vodDBName)
	vod.db = GetDB(dbName)
	// 如果找不到带 app 后缀的数据库配置，回退到默认数据库配置
	if vod.db == nil && dbName != vodDBName {
		vod.db = GetDB(vodDBName)
	}
	return vod
}

// GetList 获取视频列表
func (v *Vod) GetList(page, limit int, tag string, hour string, ids string, word string) ([]entity.VodEntity, int64, error) {
	if v.db == nil {
		return nil, 0, ErrDBConfNotFound
	}

	var vodList []entity.VodEntity
	var total int64

	db := v.db.Model(&entity.VodEntity{})

	// 按类型筛选
	if tag != "" {
		db = db.Where("type_id = ?", tag)
	}

	// 按时间筛选（最近N小时）
	if hour != "" {
		// 这里可以根据需要实现时间筛选逻辑
		// 暂时先不实现，因为需要知道具体的时间字段
	}

	// 按ID列表筛选
	if ids != "" {
		// 解析逗号分隔的ID列表
		parts := strings.Split(ids, ",")
		idList := make([]int64, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if id, err := strconv.ParseInt(part, 10, 64); err == nil {
				idList = append(idList, id)
			}
		}
		if len(idList) > 0 {
			db = db.Where("vod_id IN (?)", idList)
		}
	}

	// 按关键词搜索
	if word != "" {
		db = db.Where("vod_name LIKE ? OR vod_en LIKE ?", "%"+word+"%", "%"+word+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Order("vod_time DESC").Find(&vodList).Error; err != nil {
		return nil, 0, err
	}

	return vodList, total, nil
}

// GetByID 根据ID获取视频详情
func (v *Vod) GetByID(vodID int64) (*entity.VodEntity, error) {
	if v.db == nil {
		return nil, ErrDBConfNotFound
	}

	var vod entity.VodEntity
	err := v.db.Where("vod_id = ?", vodID).First(&vod).Error
	if err != nil {
		return nil, err
	}
	return &vod, nil
}

// Save 保存视频信息（新建或更新）
func (v *Vod) Save(vod *entity.VodEntity) error {
	if v.db == nil {
		return ErrDBConfNotFound
	}
	return v.db.Save(vod).Error
}

// BatchSave 批量保存视频信息（新建或更新）
func (v *Vod) BatchSave(vodList []*entity.VodEntity) error {
	if v.db == nil {
		return ErrDBConfNotFound
	}
	if len(vodList) == 0 {
		return ErrInvalidParam
	}
	// 使用事务批量保存
	return v.db.Transaction(func(tx *gorm.DB) error {
		for _, vod := range vodList {
			if err := tx.Save(vod).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
