package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"gitlab.com/cinemae/cine_stream/app/dao"
	"gitlab.com/cinemae/cine_stream/app/entity"
)

// ProvideService 资源提供服务
type ProvideService struct {
	ctx     context.Context
	vod     *dao.Vod
	typeDao *dao.Type
}

// NewProvideService 创建资源提供服务
func NewProvideService(ctx context.Context) *ProvideService {
	return &ProvideService{
		ctx:     ctx,
		vod:     dao.NewVod(ctx),
		typeDao: dao.NewType(ctx),
	}
}

// GetSimpleVideoList 获取简化视频列表
func (s *ProvideService) GetSimpleVideoList(page, limit int, tag, hour, ids, word string) (*entity.SimpleVideoListResponse, error) {
	vodList, total, err := s.vod.GetList(page, limit, tag, hour, ids, word)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	pageCount := int(math.Ceil(float64(total) / float64(limit)))
	if pageCount == 0 {
		pageCount = 1
	}

	// 转换为简化格式
	list := make([]entity.SimpleVideoItem, 0, len(vodList))
	for _, vod := range vodList {
		item := entity.SimpleVideoItem{
			VodID:       vod.VodID,
			VodName:     vod.VodName,
			TypeID:      vod.TypeID,
			VodEn:       vod.VodEn,
			VodRemarks:  vod.VodRemarks,
			VodPlayFrom: vod.VodPlayFrom,
			VodPlayURL:  vod.VodPlayURL,
		}

		// 格式化时间
		if vod.VodTime > 0 {
			item.VodTime = time.Unix(vod.VodTime, 0).Format("2006-01-02 15:04:05")
		}

		list = append(list, item)
	}

	// 获取分类列表
	typeList, err := s.typeDao.GetAll()
	class := make([]entity.SimpleTypeItem, 0)
	if err == nil {
		// 转换为简化格式，只包含 type_id、type_pid、type_name
		for _, t := range typeList {
			class = append(class, entity.SimpleTypeItem{
				TypeID:   t.TypeID,
				TypePID:  t.TypePID,
				TypeName: t.TypeName,
			})
		}
	}

	return &entity.SimpleVideoListResponse{
		Code:      1,
		Msg:       "数据列表",
		Page:      page,
		PageCount: pageCount,
		Limit:     strconv.Itoa(limit),
		Total:     total,
		List:      list,
		Class:     class,
	}, nil
}

// GetFullVideoList 获取完整视频列表
func (s *ProvideService) GetFullVideoList(page, limit int, tag, hour, ids, word string) (*entity.FullVideoListResponse, error) {
	vodList, total, err := s.vod.GetList(page, limit, tag, hour, ids, word)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	pageCount := int(math.Ceil(float64(total) / float64(limit)))
	if pageCount == 0 {
		pageCount = 1
	}

	// 转换为完整格式
	list := make([]entity.FullVideoItem, 0, len(vodList))
	for _, vod := range vodList {
		item := entity.FullVideoItem{
			VodEntity: vod,
		}

		// 格式化时间字段
		if vod.VodTime > 0 {
			item.VodTimeStr = time.Unix(vod.VodTime, 0).Format("2006-01-02 15:04:05")
		}

		// 格式化评分字段
		if vod.VodScore != nil {
			item.VodScoreStr = fmt.Sprintf("%.1f", *vod.VodScore)
		} else {
			item.VodScoreStr = "0.0"
		}

		// 格式化豆瓣评分字段
		if vod.VodDoubanScore != nil {
			item.VodDoubanStr = fmt.Sprintf("%.1f", *vod.VodDoubanScore)
		} else {
			item.VodDoubanStr = "0.0"
		}

		// TODO: 获取类型名称，这里暂时留空，需要查询类型表
		// item.TypeName = ...

		list = append(list, item)
	}

	return &entity.FullVideoListResponse{
		Code:      1,
		Msg:       "数据列表",
		Page:      page,
		PageCount: pageCount,
		Limit:     strconv.Itoa(limit),
		Total:     total,
		List:      list,
	}, nil
}

// Save 保存视频信息（新建或更新）
func (s *ProvideService) Save(vod *entity.VodEntity) error {
	return s.vod.Save(vod)
}

// BatchSave 批量保存视频信息（新建或更新）
func (s *ProvideService) BatchSave(vodList []*entity.VodEntity) error {
	return s.vod.BatchSave(vodList)
}
