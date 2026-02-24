// Package dal 数据层
package dao

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/cinemae/cine_stream/config"
	"gitlab.com/cinemae/cine_stream/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 错误常量定义
var (
	ErrDBConfNotFound = errors.New("数据库配置未找到")
	ErrInvalidParam   = errors.New("参数错误")
	ErrRecordNotFound = errors.New("记录不存在")
	ErrRecordExists   = errors.New("记录已存在")
)

var dbs = make(map[string]*gorm.DB)

// register 注册一个 db
func register(dbName string, db *gorm.DB) {
	dbs[dbName] = db
}

// GetDB 获取一个 db
func GetDB(dbName string) *gorm.DB {
	return dbs[dbName]
}

// InitDB 初始化数据库实例
func InitDB() {
	dbConfig := config.GetAppConf().GetDatabaseConf()
	if len(dbConfig) == 0 {
		logger.Fatalf("[InitDB] db config is empty")
	}
	for dbName, dbConf := range dbConfig {
		gormDB, err := getDB(&dbConf)
		if err != nil {
			logger.Fatalf("[InitDB] dbName=%s conn err=%+v", dbName, err)
		}
		if gormDB == nil {
			continue
		}
		register(dbName, gormDB)
	}
}

func getDB(dbConf *config.DatabaseConf) (gormDB *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConf.User,
		dbConf.Pass,
		dbConf.Host,
		dbConf.Port,
		dbConf.Name,
	)
	logger.Infof("[InitDB] db dsn=%s", dsn)

	// 配置 GORM Logger，统一打印执行的 SQL
	// 根据配置决定日志级别：如果 LogSQL 为 true，使用 Info 级别；否则使用 Silent 级别（不输出 SQL）
	logLevel := gormlogger.Silent
	if dbConf.LogSQL {
		logLevel = gormlogger.Info
	}
	gormLogger := gormlogger.New(
		&gormLogWriter{logSQL: dbConf.LogSQL},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢查询阈值
			LogLevel:                  logLevel,               // 日志级别，根据 LogSQL 配置决定
			IgnoreRecordNotFoundError: true,                   // 忽略 ErrRecordNotFound 错误
			Colorful:                  false,                  // 禁用彩色打印，使用统一的日志格式
		},
	)

	gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   dbConf.TablePrefix,
			SingularTable: false,
		},
		Logger: gormLogger,
	})
	if err != nil {
		return gormDB, err
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return gormDB, err
	}
	if dbConf.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(dbConf.ConnMaxIdle)
	}
	if dbConf.ConnMaxConnection > 0 {
		sqlDB.SetMaxOpenConns(dbConf.ConnMaxConnection)
	}
	if dbConf.ConnMaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifeTime) * time.Second)
	}
	return gormDB, nil
}

// CloseDB 关闭单个 db
func CloseDB(dbName string) {
	gormDB := GetDB(dbName)
	if gormDB == nil {
		return
	}
	db, err := gormDB.DB()
	if err != nil {
		logger.Errorf("[CloseDB] dbName=%s err: %v", dbName, err)
		return
	}
	db.Close()
}

// CloseDBs 关闭所有的 db
func CloseDBs() {
	for dbName, _ := range dbs {
		CloseDB(dbName)
	}
}

// gormLogWriter 实现 GORM logger.Writer 接口，使用项目的日志系统
type gormLogWriter struct {
	logSQL bool // 是否输出 SQL 日志
}

// Printf 实现 logger.Writer 接口
func (w *gormLogWriter) Printf(format string, v ...interface{}) {
	// 只有当 logSQL 为 true 时才输出 SQL 日志
	if w.logSQL {
		logger.Infof("[GORM] "+format, v...)
	}
}
