// Package service 业务逻辑层
package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/config"
	"gitlab.com/cinemae/cine_stream/logger"
)

// PassportPlayRightsResponse Passport 播放权限响应结构
type PassportPlayRightsResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Sub    string `json:"sub"`
	Name   string `json:"name"`
	Data   struct {
		Definition    string `json:"definition"`
		HasPermission bool   `json:"hasPermission"`
		Owner         string `json:"owner"`
		Reason        string `json:"reason"`
		UserId        string `json:"userId"`
		VideoId       string `json:"videoId"`
		VodId         string `json:"vodId"`
	} `json:"data"`
	Data2 interface{} `json:"data2"`
	Data3 interface{} `json:"data3"`
}

// CheckPlayRights 检查用户是否有播放权限
// 根据 access_token 去 passport 获取当前用户是否是播放权限
// 有播放权限返回 true，没有播放权限返回 false
func CheckPlayRights(ctx *gin.Context, videoID string) bool {
	// 获取 access_token
	accessToken := entity.ContextValueLoginToken(ctx)
	if accessToken == "" {
		logger.WithContext(ctx).Warnf("[CheckPlayRights] access_token 为空")
		return false
	}

	// 获取 passport 配置（包含播放权限接口配置）
	passportConf := config.GetAppConf().GetPassportConf()
	if passportConf.Endpoint == "" {
		logger.WithContext(ctx).Errorf("[CheckPlayRights] passport endpoint 配置为空")
		return false
	}

	// 获取播放权限接口路径，默认为 /api/get-user-play-rights
	playRightsAPI := passportConf.PlayRightsAPI
	if playRightsAPI == "" {
		playRightsAPI = "/api/get-user-play-rights"
	}

	// 确保路径以 / 开头
	if !strings.HasPrefix(playRightsAPI, "/") {
		playRightsAPI = "/" + playRightsAPI
	}

	// 构建完整的请求 URL
	endpoint := strings.TrimSuffix(passportConf.Endpoint, "/")
	reqURL := fmt.Sprintf("%s%s?video_id=%s", endpoint, playRightsAPI, url.QueryEscape(videoID))

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		logger.WithContext(ctx).Errorf("[CheckPlayRights] 创建请求失败: %v", err)
		return false
	}

	// 设置 Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.WithContext(ctx).Errorf("[CheckPlayRights] 请求 passport 失败: %v", err)
		return false
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithContext(ctx).Errorf("[CheckPlayRights] 读取响应失败: %v", err)
		return false
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		logger.WithContext(ctx).Warnf("[CheckPlayRights] passport 返回非 200 状态码: %d, body: %s", resp.StatusCode, string(body))
		return false
	}

	// 解析响应
	var playRightsResp PassportPlayRightsResponse
	if err := json.Unmarshal(body, &playRightsResp); err != nil {
		logger.WithContext(ctx).Errorf("[CheckPlayRights] 解析响应失败: %v, body: %s", err, string(body))
		return false
	}

	// 检查响应状态
	if playRightsResp.Status != "ok" {
		logger.WithContext(ctx).Warnf("[CheckPlayRights] passport 返回错误状态: status=%s, msg=%s", playRightsResp.Status, playRightsResp.Msg)
		return false
	}

	// 检查是否有权限
	return playRightsResp.Data.HasPermission
}
