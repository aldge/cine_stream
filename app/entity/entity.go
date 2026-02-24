// Package entity 实体定义
package entity

// Response 返回结果定义
type Response struct {
	Code    int32       `json:"code"`    // 错误码
	Message string      `json:"message"` // 返回信息
	Data    interface{} `json:"data"`    // 返回数据
}
