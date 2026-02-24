// Package cmd 命令行参数解析
package cmd

import (
	"flag"
	"os"
)

const (
	DefaultAppConfigPath = "./app.yaml"    // 默认的 配置文件地址
	DefaultAppConfigName = "app.yaml"      // 默认的配置文件名
	AppEnvVarName        = "KAPP_ENV_TYPE" // 环境变量名
)

var (
	FlagVar = &flagVar{} // FlagVar 命令行参数
	EnvVar  = &envVar{}  // EnvVar 环境变量
)

type flagVar struct {
	runEnv      string // 运行环境
	appConfPath string // 项目配置文件地址
	migrate     bool   // 是否执行迁移
	migrateDB   string // 指定要迁移的数据库名称，为空则迁移所有数据库
}

type envVar struct {
	runEnv string // 运行环境
}

// GetAppConfPath 获取 app 配置路径
func (fv *flagVar) GetAppConfPath() string {
	return FlagVar.appConfPath
}

// GetRunEnv 获取环境
func GetRunEnv() string {
	if FlagVar.runEnv != "" {
		return FlagVar.runEnv
	}
	return EnvVar.runEnv
}

// initFlagVar 初始化解析命令行参数
func initFlagVar() {
	flag.StringVar(&FlagVar.runEnv, "env", "", "env, value can be debug/dev/test/prod")
	flag.StringVar(&FlagVar.appConfPath, "conf", DefaultAppConfigPath, "server config path")
	flag.BoolVar(&FlagVar.migrate, "migrate", false, "执行数据库迁移")
	flag.StringVar(&FlagVar.migrateDB, "migrate-db", "", "指定要迁移的数据库名称，为空则迁移所有数据库")
	flag.Parse()
}

// GetMigrate 获取是否执行迁移
func (fv *flagVar) GetMigrate() bool {
	return FlagVar.migrate
}

// GetMigrateDB 获取要迁移的数据库名称
func (fv *flagVar) GetMigrateDB() string {
	return FlagVar.migrateDB
}

// initEnvVar 初始化环境变量
func initEnvVar() {
	// 获取环境变量中配置的
	EnvVar.runEnv = os.Getenv(AppEnvVarName)
}

// InitFlagEnvVar 初始化命令行和环境变量
func Init() {
	initFlagVar()
	initEnvVar()
}
