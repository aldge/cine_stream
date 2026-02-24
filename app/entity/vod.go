package entity

import (
	"encoding/json"
)

// TypeEntity 影视类型实体
// 对应数据库表 cine_type
// 详细字段说明请参考 migrations/cine_type-20260104-01.sql
type TypeEntity struct {
	TypeID     uint16  `gorm:"column:type_id;primaryKey;autoIncrement" json:"type_id"`
	TypeName   *string `gorm:"column:type_name;size:60" json:"type_name"`
	TypeEn     *string `gorm:"column:type_en;size:60" json:"type_en"`
	TypePID    uint16  `gorm:"column:type_pid;default:0" json:"type_pid"`
	TypeStatus uint8   `gorm:"column:type_status;default:1" json:"type_status"`
}

// TableName 指定表名
func (TypeEntity) TableName() string {
	return "cine_type"
}

// VodEntity VOD实体
// 对应数据库表 cine_vod
// 详细字段说明请参考 migrations/cine_vod-20250104-01.sql
type VodEntity struct {
	VodID            int64    `gorm:"column:vod_id;primaryKey;autoIncrement" json:"vod_id"`
	TypeID           *int64   `gorm:"column:type_id" json:"type_id"`
	TypeID1          *int64   `gorm:"column:type_id_1" json:"type_id_1"`
	GroupID          *int64   `gorm:"column:group_id" json:"group_id"`
	VodName          *string  `gorm:"column:vod_name;size:255" json:"vod_name"`
	VodSub           *string  `gorm:"column:vod_sub;size:255" json:"vod_sub"`
	VodEn            *string  `gorm:"column:vod_en;size:255" json:"vod_en"`
	VodStatus        *int8    `gorm:"column:vod_status" json:"vod_status"`
	VodLetter        *string  `gorm:"column:vod_letter;size:191" json:"vod_letter"`
	VodColor         *string  `gorm:"column:vod_color;size:6" json:"vod_color"`
	VodTag           *string  `gorm:"column:vod_tag;size:100" json:"vod_tag"`
	VodClass         *string  `gorm:"column:vod_class;size:255" json:"vod_class"`
	VodPic           *string  `gorm:"column:vod_pic;size:1024" json:"vod_pic"`
	VodPicThumb      *string  `gorm:"column:vod_pic_thumb;size:1024" json:"vod_pic_thumb"`
	VodPicSlide      *string  `gorm:"column:vod_pic_slide;size:1024" json:"vod_pic_slide"`
	VodPicScreenshot *string  `gorm:"column:vod_pic_screenshot;size:191" json:"vod_pic_screenshot"`
	VodActor         *string  `gorm:"column:vod_actor;size:255" json:"vod_actor"`
	VodDirector      *string  `gorm:"column:vod_director;size:255" json:"vod_director"`
	VodWriter        *string  `gorm:"column:vod_writer;size:100" json:"vod_writer"`
	VodBehind        *string  `gorm:"column:vod_behind;size:100" json:"vod_behind"`
	VodBlurb         *string  `gorm:"column:vod_blurb;type:text" json:"vod_blurb"`
	VodRemarks       *string  `gorm:"column:vod_remarks;size:100" json:"vod_remarks"`
	VodPubdate       *string  `gorm:"column:vod_pubdate;size:100" json:"vod_pubdate"`
	VodTotal         int64    `gorm:"column:vod_total;default:1" json:"vod_total"`
	VodSerial        *string  `gorm:"column:vod_serial;size:20" json:"vod_serial"`
	VodTV            *string  `gorm:"column:vod_tv;size:30" json:"vod_tv"`
	VodWeekday       *string  `gorm:"column:vod_weekday;size:30" json:"vod_weekday"`
	VodArea          *string  `gorm:"column:vod_area;size:20" json:"vod_area"`
	VodLang          *string  `gorm:"column:vod_lang;size:10" json:"vod_lang"`
	VodYear          *string  `gorm:"column:vod_year;size:10" json:"vod_year"`
	VodVersion       *string  `gorm:"column:vod_version;size:30" json:"vod_version"`
	VodState         *string  `gorm:"column:vod_state;size:30" json:"vod_state"`
	VodAuthor        *string  `gorm:"column:vod_author;size:60" json:"vod_author"`
	VodJumpurl       *string  `gorm:"column:vod_jumpurl;size:150" json:"vod_jumpurl"`
	VodTpl           *string  `gorm:"column:vod_tpl;size:30" json:"vod_tpl"`
	VodTplPlay       *string  `gorm:"column:vod_tpl_play;size:30" json:"vod_tpl_play"`
	VodTplDown       *string  `gorm:"column:vod_tpl_down;size:30" json:"vod_tpl_down"`
	VodIsend         *int8    `gorm:"column:vod_isend" json:"vod_isend"`
	VodLock          *int8    `gorm:"column:vod_lock" json:"vod_lock"`
	VodLevel         *int8    `gorm:"column:vod_level" json:"vod_level"`
	VodCopyright     *int8    `gorm:"column:vod_copyright" json:"vod_copyright"`
	VodPoints        int64    `gorm:"column:vod_points;default:0" json:"vod_points"`
	VodPointsPlay    int64    `gorm:"column:vod_points_play;default:0" json:"vod_points_play"`
	VodPointsDown    int64    `gorm:"column:vod_points_down;default:0" json:"vod_points_down"`
	VodHits          int64    `gorm:"column:vod_hits;default:0" json:"vod_hits"`
	VodHitsDay       int64    `gorm:"column:vod_hits_day;default:0" json:"vod_hits_day"`
	VodHitsWeek      int64    `gorm:"column:vod_hits_week;default:0" json:"vod_hits_week"`
	VodHitsMonth     int64    `gorm:"column:vod_hits_month;default:0" json:"vod_hits_month"`
	VodDuration      *string  `gorm:"column:vod_duration;size:10" json:"vod_duration"`
	VodUp            int64    `gorm:"column:vod_up;default:0" json:"vod_up"`
	VodDown          int64    `gorm:"column:vod_down;default:0" json:"vod_down"`
	VodScore         *float32 `gorm:"column:vod_score" json:"vod_score"`
	VodScoreAll      int64    `gorm:"column:vod_score_all;default:0" json:"vod_score_all"`
	VodScoreNum      int64    `gorm:"column:vod_score_num;default:0" json:"vod_score_num"`
	VodTime          int64    `gorm:"column:vod_time;default:0" json:"vod_time"`
	VodTimeAdd       int64    `gorm:"column:vod_time_add;default:0" json:"vod_time_add"`
	VodTimeHits      int64    `gorm:"column:vod_time_hits;default:0" json:"vod_time_hits"`
	VodTimeMake      int64    `gorm:"column:vod_time_make;default:0" json:"vod_time_make"`
	VodTrysee        int64    `gorm:"column:vod_trysee;default:0" json:"vod_trysee"`
	VodDoubanID      int64    `gorm:"column:vod_douban_id;default:0" json:"vod_douban_id"`
	VodDoubanScore   *float32 `gorm:"column:vod_douban_score" json:"vod_douban_score"`
	VodReurl         *string  `gorm:"column:vod_reurl;size:255" json:"vod_reurl"`
	VodRelVod        *string  `gorm:"column:vod_rel_vod;size:255" json:"vod_rel_vod"`
	VodRelArt        *string  `gorm:"column:vod_rel_art;size:255" json:"vod_rel_art"`
	VodPwd           *string  `gorm:"column:vod_pwd;size:10" json:"vod_pwd"`
	VodPwdURL        *string  `gorm:"column:vod_pwd_url;size:255" json:"vod_pwd_url"`
	VodPwdPlay       *string  `gorm:"column:vod_pwd_play;size:10" json:"vod_pwd_play"`
	VodPwdPlayURL    *string  `gorm:"column:vod_pwd_play_url;size:255" json:"vod_pwd_play_url"`
	VodPwdDown       *string  `gorm:"column:vod_pwd_down;size:10" json:"vod_pwd_down"`
	VodPwdDownURL    *string  `gorm:"column:vod_pwd_down_url;size:255" json:"vod_pwd_down_url"`
	VodContent       *string  `gorm:"column:vod_content;type:text" json:"vod_content"`
	VodPlayFrom      *string  `gorm:"column:vod_play_from;size:255" json:"vod_play_from"`
	VodPlayServer    *string  `gorm:"column:vod_play_server;size:255" json:"vod_play_server"`
	VodPlayNote      *string  `gorm:"column:vod_play_note;size:255" json:"vod_play_note"`
	VodPlayURL       *string  `gorm:"column:vod_play_url;type:text" json:"vod_play_url"`
	VodDownFrom      *string  `gorm:"column:vod_down_from;size:255" json:"vod_down_from"`
	VodDownServer    *string  `gorm:"column:vod_down_server;size:255" json:"vod_down_server"`
	VodDownNote      *string  `gorm:"column:vod_down_note;size:255" json:"vod_down_note"`
	VodDownURL       *string  `gorm:"column:vod_down_url;size:191" json:"vod_down_url"`
	VodPlot          *int8    `gorm:"column:vod_plot" json:"vod_plot"`
	VodPlotName      *string  `gorm:"column:vod_plot_name;size:191" json:"vod_plot_name"`
	VodPlotDetail    *string  `gorm:"column:vod_plot_detail;type:text" json:"vod_plot_detail"`
}

// TableName 指定表名
func (VodEntity) TableName() string {
	return "cine_vod"
}

// UnmarshalJSON 自定义JSON反序列化，处理布尔值到int8的转换
func (v *VodEntity) UnmarshalJSON(data []byte) error {
	// 使用临时结构体来处理布尔值字段
	type Alias VodEntity
	aux := &struct {
		VodIsend     interface{} `json:"vod_isend"`
		VodLock      interface{} `json:"vod_lock"`
		VodLevel     interface{} `json:"vod_level"`
		VodCopyright interface{} `json:"vod_copyright"`
		VodPlot      interface{} `json:"vod_plot"`
		VodStatus    interface{} `json:"vod_status"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 辅助函数：将interface{}转换为*int8
	convertToInt8Ptr := func(val interface{}) *int8 {
		if val == nil {
			return nil
		}
		if b, ok := val.(bool); ok {
			if b {
				result := int8(1)
				return &result
			}
			result := int8(0)
			return &result
		}
		if num, ok := val.(float64); ok {
			result := int8(num)
			return &result
		}
		return nil
	}

	// 转换布尔值到int8
	v.VodIsend = convertToInt8Ptr(aux.VodIsend)
	v.VodLock = convertToInt8Ptr(aux.VodLock)
	v.VodLevel = convertToInt8Ptr(aux.VodLevel)
	v.VodCopyright = convertToInt8Ptr(aux.VodCopyright)
	v.VodPlot = convertToInt8Ptr(aux.VodPlot)
	v.VodStatus = convertToInt8Ptr(aux.VodStatus)

	return nil
}

// SimpleVideoListResponse 简化视频列表响应
type SimpleVideoListResponse struct {
	Code      int               `json:"code"`
	Msg       string            `json:"msg"`
	Page      int               `json:"page"`
	PageCount int               `json:"pagecount"`
	Limit     string            `json:"limit"`
	Total     int64             `json:"total"`
	List      []SimpleVideoItem `json:"list"`
	Class     []SimpleTypeItem  `json:"class"`
}

// SimpleTypeItem 简化类型项（只包含 type_id、type_pid、type_name）
type SimpleTypeItem struct {
	TypeID   uint16  `json:"type_id"`
	TypePID  uint16  `json:"type_pid"`
	TypeName *string `json:"type_name"`
}

// SimpleVideoItem 简化视频项
type SimpleVideoItem struct {
	VodID       int64   `json:"vod_id"`
	VodName     *string `json:"vod_name"`
	TypeID      *int64  `json:"type_id"`
	TypeName    *string `json:"type_name"`
	VodEn       *string `json:"vod_en"`
	VodTime     string  `json:"vod_time"`
	VodRemarks  *string `json:"vod_remarks"`
	VodPlayFrom *string `json:"vod_play_from"`
	VodPlayURL  *string `json:"vod_play_url"`
}

// FullVideoListResponse 完整视频列表响应
type FullVideoListResponse struct {
	Code      int             `json:"code"`
	Msg       string          `json:"msg"`
	Page      int             `json:"page"`
	PageCount int             `json:"pagecount"`
	Limit     string          `json:"limit"`
	Total     int64           `json:"total"`
	List      []FullVideoItem `json:"list"`
}

// FullVideoItem 完整视频项
type FullVideoItem struct {
	VodEntity
	TypeName     string `json:"type_name"`
	VodTimeStr   string `json:"vod_time"`         // 格式化的时间字符串（覆盖原始字段）
	VodScoreStr  string `json:"vod_score"`        // 格式化的评分字符串（覆盖原始字段）
	VodDoubanStr string `json:"vod_douban_score"` // 格式化的豆瓣评分字符串（覆盖原始字段）
}

// MarshalJSON 自定义JSON序列化，确保时间、评分字段格式正确
func (f FullVideoItem) MarshalJSON() ([]byte, error) {
	// 创建一个map来存储所有字段
	m := make(map[string]interface{})

	// 先序列化嵌入的VodEntity
	vodBytes, err := json.Marshal(f.VodEntity)
	if err != nil {
		return nil, err
	}
	var vodMap map[string]interface{}
	if err := json.Unmarshal(vodBytes, &vodMap); err != nil {
		return nil, err
	}

	// 复制所有字段到map
	for k, v := range vodMap {
		m[k] = v
	}

	// 覆盖时间、评分字段为字符串格式
	if f.VodTimeStr != "" {
		m["vod_time"] = f.VodTimeStr
	}
	if f.VodScoreStr != "" {
		m["vod_score"] = f.VodScoreStr
	}
	if f.VodDoubanStr != "" {
		m["vod_douban_score"] = f.VodDoubanStr
	}

	// 添加类型名称
	if f.TypeName != "" {
		m["type_name"] = f.TypeName
	}

	return json.Marshal(m)
}
