package dao

import (
	"context"
	"hash/fnv"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/cinemae/gopkg/app"
)

// Base 基础的 dao
type Base struct {
}

// hashShardingID 计算分表ID的哈希值（返回uint64，避免负数）
func hashShardingID(shardingID string) uint64 {
	if shardingID == "" {
		return 0 // 空ID映射到默认表
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(shardingID))
	return h.Sum64()
}

// getShardingTableName 获取分表的表名，固定格式 xxx_0 xxx_1
func getShardingTableName(shardingID string, tableNum int, tableNamePrefix string) string {
	if tableNum <= 1 {
		return tableNamePrefix
	}
	hashValue := hashShardingID(shardingID)
	tableIndex := hashValue % uint64(tableNum)
	return tableNamePrefix + "_" + strconv.Itoa(int(tableIndex))
}

// getAppDBName 根据不同的app获取不同的db；格式：{defaultDbName}_{appName}
func getAppDBName(ctx context.Context, defaultDbName string) string {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return defaultDbName
	}
	appName := app.GetAppName(ginCtx)
	if appName == "" {
		return defaultDbName
	}
	return defaultDbName + "_" + string(appName)
}
