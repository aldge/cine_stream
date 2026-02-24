package entity

// VideoEncryptEntity 视频加密信息实体
// 对应数据库表 cine_video_encrypt
// 详细字段说明请参考 docs/video.sql
type VideoEncryptEntity struct {
	VideoEncryptID uint64 `gorm:"column:video_encrypt_id;primaryKey;autoIncrement" json:"video_encrypt_id"`
	VideoID        string `gorm:"column:video_id;size:32;not null;index" json:"video_id"`
	Key            string `gorm:"column:key;size:64;not null" json:"key"`
	IV             string `gorm:"column:iv;size:64;not null" json:"iv"`
	CreateTime     uint64 `gorm:"column:create_time;not null;index" json:"create_time"`
}
