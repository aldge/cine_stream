package main

import (
	"net/http"
	"os"
	"runtime"

	"gitlab.com/cinemae/cine_stream/app/dao"
	"gitlab.com/cinemae/cine_stream/cmd"
	"gitlab.com/cinemae/cine_stream/config"
	"gitlab.com/cinemae/cine_stream/logger"
	"gitlab.com/cinemae/cine_stream/router"
	"gitlab.com/cinemae/cine_stream/utils"

	"github.com/gin-gonic/gin"
	passport "gitlab.com/cinemae/gopkg/casdoor"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 命令行初始化
	cmd.Init()
	// 配置初始化
	config.Init()
	// 日志初始化
	logger.Init()
	defer logger.Sync()
	// 初始化 Passport SDK (使用 Casdoor 开源项目)
	initPassport()
	// 路由和中间件初始化
	router.Init()
	// DB 初始化
	dao.InitDB()

	// 如果指定了迁移参数，执行迁移并退出
	if cmd.FlagVar.GetMigrate() {
		dbName := cmd.FlagVar.GetMigrateDB()
		// 获取数据库配置并转换为迁移所需的格式
		dbConfigMap := config.GetAppConf().GetDatabaseConf()
		migrateDbConfig := make(map[string]utils.DatabaseConf)
		for name, dbConf := range dbConfigMap {
			migrateDbConfig[name] = utils.DatabaseConf{
				Host: dbConf.Host,
				Port: dbConf.Port,
				Name: dbConf.Name,
				User: dbConf.User,
				Pass: dbConf.Pass,
			}
		}
		if err := utils.RunMigrations(migrateDbConfig, dbName); err != nil {
			logger.Fatalf("[main] 执行迁移失败: %v", err)
		}
		logger.Infof("[main] 迁移执行完成")
		os.Exit(0)
	}

	// 设置 gin 框架允许环境
	gin.SetMode(config.GetAppConf().Global.GinMode)

	// 启动 server
	s := &http.Server{
		Addr:           config.GetServerAddr(),
		Handler:        router.GinEngine,
		ReadTimeout:    config.GetReadTimeout(),
		WriteTimeout:   config.GetWriteTimeout(),
		MaxHeaderBytes: 1 << 20,
	}
	logger.Infof("[cine_server] Listen addr=%+v", config.GetServerAddr())
	err := s.ListenAndServe()
	if err != nil {
		logger.Fatalf("[cine_server] err:%s", err)
	}
}

// initPassport 初始化 Passport SDK (使用 Casdoor 开源项目)
func initPassport() {
	passportConf := config.GetAppConf().GetPassportConf()
	if passportConf.Endpoint == "" || passportConf.ClientID == "" {
		logger.Warnf("[InitPassport] Passport config is not complete, skip initialization")
		return
	}

	// 读取证书文件内容
	var certificate string
	if passportConf.CertificatePath != "" {
		certBytes, err := os.ReadFile(passportConf.CertificatePath)
		if err != nil {
			logger.Fatalf("[InitPassport] Failed to read certificate file: %v", err)
			return
		}
		certificate = string(certBytes)
	}

	// 初始化 Casdoor SDK (Passport 使用 Casdoor 开源项目)
	passport.InitConfig(
		passportConf.Endpoint,
		passportConf.ClientID,
		passportConf.ClientSecret,
		certificate,
		passportConf.OrganizationName,
		passportConf.ApplicationName,
	)
	logger.Infof("[InitPassport] Passport SDK initialized successfully")
}
