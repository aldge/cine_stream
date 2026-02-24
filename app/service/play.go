package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/cinemae/cine_stream/app/dao"
	"gitlab.com/cinemae/cine_stream/app/entity"
	"gitlab.com/cinemae/cine_stream/config"
	"gitlab.com/cinemae/cine_stream/utils"
	"gitlab.com/cinemae/gopkg/app"
)

// Play 播放业务逻辑
type Play struct {
	ctx             context.Context
	daoVideoEncrypt *dao.VideoEncrypt
}

// NewPlay 创建TS切片业务逻辑对象
func NewPlay(ctx context.Context) *Play {
	return &Play{
		ctx:             ctx,
		daoVideoEncrypt: dao.NewVideoEncrypt(ctx),
	}
}

// GenerateM3U8Content 生成M3U8文件内容（支持每个切片独立的加密信息）
func (p *Play) GenerateM3U8Content(ctx *gin.Context, videoID string, tsList []entity.VideoTSEntity) (string, error) {
	if videoID == "" {
		return "", errors.New("视频ID不能为空")
	}
	if len(tsList) == 0 {
		return "", errors.New("该视频没有TS切片")
	}

	// 获取所有切片的加密信息
	encryptInfo, err := p.daoVideoEncrypt.GetByVideoID(videoID)
	if err != nil {
		return "", errors.New("获取视频加密信息失败")
	}

	baseURL := utils.GetRequestBaseURL(ctx)
	appName := app.GetAppName(ctx)

	// 计算最大时长（TARGETDURATION 应该是所有片段的最大时长，向上取整）
	maxDuration := 0.0
	for _, ts := range tsList {
		if ts.Duration > maxDuration {
			maxDuration = ts.Duration
		}
	}
	targetDuration := int(math.Ceil(maxDuration))
	if targetDuration < 1 {
		targetDuration = 1 // 最小值设为 1
	}

	// 生成M3U8文件内容（按照标准顺序）
	m3u8Content := "#EXTM3U\n"
	m3u8Content += "#EXT-X-VERSION:3\n"
	m3u8Content += "#EXT-X-MEDIA-SEQUENCE:0\n"
	m3u8Content += "#EXT-X-ALLOW-CACHE:YES\n"
	m3u8Content += fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", targetDuration)
	m3u8Content += fmt.Sprintf(`#EXT-X-KEY:METHOD=AES-128,URI="%s/play/key/%s?app=%s",IV=0x%s`+"\n", baseURL, videoID, appName, encryptInfo.IV)
	// 为每个切片添加信息（#EXTINF ）
	for _, ts := range tsList {
		m3u8Content += "#EXTINF:" + formatDuration(ts.Duration) + ",\n"
		m3u8Content += buildTsUrl(ts.TSPath) + "\n"
	}
	m3u8Content += "#EXT-X-ENDLIST\n"
	return m3u8Content, nil
}

func buildTsUrl(tsPath string) string {
	// 判断 tsPath 是否包含域名
	if strings.HasPrefix(tsPath, "http://") || strings.HasPrefix(tsPath, "https://") {
		return tsPath
	}
	// 获取 cdn 域名
	cdnConf := config.GetAppConf().GetCDNConf()
	if len(cdnConf) == 0 {
		return tsPath
	}
	// 如果 tsPath 带 / ，去掉
	tsPath = strings.TrimPrefix(tsPath, "/")
	// todo 这里可以根据不同的地域返回不同的 cdn 域名，暂时先使用默认
	defaultCdnConf := cdnConf["default"]
	return fmt.Sprintf("%s/%s", defaultCdnConf.URL, tsPath)
}

// formatDuration 格式化时长
func formatDuration(duration float64) string {
	return strconv.FormatFloat(duration, 'f', 6, 64)
}
