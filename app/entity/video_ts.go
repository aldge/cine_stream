package entity

// VideoTSEntity TS切片实体
// 对应数据库表 cine_video_ts
// 详细字段说明请参考 docs/video.sql
type VideoTSEntity struct {
	VideoTSID  int64   `gorm:"column:video_ts_id;primaryKey;autoIncrement" json:"video_ts_id"`
	VideoID    string  `gorm:"column:video_id" json:"video_id"`
	TSSequence int64   `gorm:"column:ts_sequence" json:"ts_sequence"`
	TSPath     string  `gorm:"column:ts_path" json:"ts_path"`
	Duration   float64 `gorm:"column:duration" json:"duration"`
	Definition string  `gorm:"column:definition" json:"definition"`
	CreateTime int64   `gorm:"column:create_time" json:"create_time"`
}

// VideoTSSaveRequest 批量保存TS切片请求参数
type VideoTSSaveRequest struct {
	VideoID string                 `json:"video_id" binding:"required"`
	Key     string                 `json:"key" binding:"required"`
	IV      string                 `json:"iv" binding:"required"`
	TSData  []*VideoTsSaveDataItem `json:"ts_data" binding:"required"`
}

// VideoTsSaveDataItem 批量保存TS切片请求参数中的单个TS切片数据
type VideoTsSaveDataItem struct {
	TSSequence int64   `json:"ts_sequence" binding:"required"`
	TSPath     string  `json:"ts_path" binding:"required"`
	Duration   float64 `json:"duration" binding:"required"`
	Definition string  `json:"definition"`
}
