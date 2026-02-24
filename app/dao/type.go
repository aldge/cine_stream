package dao

import (
	"context"

	"gitlab.com/cinemae/cine_stream/app/entity"
	"gorm.io/gorm"
)

const (
	typeDBName = "cine_stream" // Type 表数据库名
)

// Type 类型数据访问对象
type Type struct {
	ctx context.Context
	db  *gorm.DB
}

// NewType 创建类型数据访问对象
func NewType(ctx context.Context) *Type {
	t := &Type{
		ctx: ctx,
	}
	dbName := getAppDBName(ctx, typeDBName)
	t.db = GetDB(dbName)
	// 如果找不到带 app 后缀的数据库配置，回退到默认数据库配置
	if t.db == nil && dbName != typeDBName {
		t.db = GetDB(typeDBName)
	}
	return t
}

// GetAll 获取所有类型列表
func (t *Type) GetAll() ([]entity.TypeEntity, error) {
	if t.db == nil {
		return nil, ErrDBConfNotFound
	}

	var typeList []entity.TypeEntity
	err := t.db.Model(&entity.TypeEntity{}).
		Where("type_status = ?", 1).
		Order("type_id ASC").
		Find(&typeList).Error
	if err != nil {
		return nil, err
	}
	return typeList, nil
}
