package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gitlab.com/cinemae/cine_stream/cmd"
	"gitlab.com/cinemae/cine_stream/utils"
	klog "gitlab.com/cinemae/gopkg/log"
	"gopkg.in/yaml.v3"
)

var (
	appConfig *AppConfig
)

// AppConfig 配置数据结构
type AppConfig struct {
	// Global 全局配置
	Global struct {
		Env     string `yaml:"env"`      // 运行环境
		GinMode string `yaml:"gin_mode"` // gin 框架运行环境
	} `yaml:"Global"`
	// Server 服务配置
	Server struct {
		IP           string `yaml:"ip"`            // ip
		Port         int    `yaml:"port"`          // port
		ReadTimeout  int    `yaml:"read_timeout"`  // 读超时时间 ms
		WriteTimeout int    `yaml:"write_timeout"` // 写超时时间 ms
	} `yaml:"Server"`
	// Database 数据库配置
	Database map[string]DatabaseConf `yaml:"Database"`
	// CDN 配置
	CDN map[string]CDNConf `yaml:"CDN"`
	// Logger 日志配置
	Logger map[string]klog.Config `yaml:"Logger"`
	// Auth 登录认证配置
	Auth AuthConf `yaml:"Auth"`
}

// DatabaseConf 数据库配置
type DatabaseConf struct {
	Host              string      `yaml:"host"`                // host
	Port              int         `yaml:"port"`                // 端口
	Name              string      `yaml:"name"`                // 数据库名
	User              string      `yaml:"user"`                // 用户名
	Pass              string      `yaml:"pass"`                // 密码
	TablePrefix       string      `yaml:"table_prefix"`        // 表前缀
	ConnMaxIdle       int         `yaml:"conn_max_idle"`       // 连接最大空闲数
	ConnMaxConnection int         `yaml:"conn_max_connection"` // 最大连接数
	ConnMaxLifeTime   int         `yaml:"conn_max_lifetime"`   // 连接最大活动时间
	TableConfig       []TableConf `yaml:"table_config"`        // 分表配置
	LogSQL            bool        `yaml:"log_sql"`             // 是否输出 SQL 日志，默认 false
}

// TableConf 数据库表配置
type TableConf struct {
	TableName   string `yaml:"table_name"`   // 表名称
	ShardingNum int    `yaml:"sharding_num"` // 分表数量
}

// AuthConf 登录配置
type AuthConf struct {
	JwtSecret   string       `yaml:"jwt_secret"`   // jwt 密匙
	ExpireHours int          `yaml:"expire_hours"` // 过期时间小时
	Passport    PassportConf `yaml:"passport"`     // Passport 配置
}

// PassportConf Passport 配置
type PassportConf struct {
	Endpoint         string `yaml:"endpoint"`          // Passport 服务地址
	ClientID         string `yaml:"client_id"`         // 客户端ID
	ClientSecret     string `yaml:"client_secret"`     // 客户端密钥
	CertificatePath  string `yaml:"certificate_path"`  // 证书文件路径
	OrganizationName string `yaml:"organization_name"` // 组织名称
	ApplicationName  string `yaml:"application_name"`  // 应用名称
	PlayRightsAPI    string `yaml:"play_rights_api"`   // 播放权限接口路径（相对于 endpoint）
}

// CDNConf CDN 配置
type CDNConf struct {
	URL string `yaml:"url"` // CDN URL
}

// getAppConfigPath 获取服务启动配置文件路径
//
//	-conf 传入配置文件路径
//	默认路径 ./app.yaml
func getAppConfigPath() (string, error) {
	if cmd.FlagVar.GetAppConfPath() != cmd.DefaultAppConfigPath {
		return cmd.FlagVar.GetAppConfPath(), nil
	}
	path, err := utils.SearchPath(cmd.DefaultAppConfigName, cmd.GetRunEnv())
	return path, err
}

// initAppConfig 初始化项目配置
func initAppConfig() {
	// 获取项目配置文件路径
	path, err := getAppConfigPath()
	if err != nil {
		panic("get app config path fail: " + err.Error())
	}
	log.Printf("[Config] get config file: %s \n", path)
	// 解析项目配置
	cfg, err := LoadAppConfig(path)
	if err != nil {
		panic("parse config fail: " + err.Error())
	}
	//log.Println(fmt.Sprintf("[Config] config: %+v", cfg))
	SetAppConf(cfg)
}

// CorrectConfig 修正配置
func CorrectConfig(config *AppConfig) error {
	return nil
}

// parseAppConfigYaml 解析配置从 yaml
func parseAppConfigYaml(configPath string) (*AppConfig, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	cfg := defaultAppConfig()
	if err := yaml.Unmarshal(buf, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// 默认的配置
func defaultAppConfig() *AppConfig {
	cfg := &AppConfig{}
	return cfg
}

// LoadAppConfig 从配置文件加载项目配置
func LoadAppConfig(configPath string) (*AppConfig, error) {
	cfg, err := parseAppConfigYaml(configPath)
	if err != nil {
		return nil, err
	}
	if err := CorrectConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// SetGlobalAppConf 设置全局的 app 配置
func SetAppConf(config *AppConfig) {
	appConfig = config
}

// GetAppConf 设置全局的 app 配置
func GetAppConf() *AppConfig {
	return appConfig
}

// GetServerAddr 获取 Server 监听的IP和端口
func GetServerAddr() string {
	return fmt.Sprintf("%s:%d", GetAppConf().Server.IP, GetAppConf().Server.Port)
}

// GetReadTimeout 读超时时间
func GetReadTimeout() time.Duration {
	if GetAppConf().Server.ReadTimeout > 0 {
		return time.Duration(GetAppConf().Server.ReadTimeout) * time.Millisecond
	}
	return 2 * time.Second
}

// GetWriteTimeout 写超时时间
func GetWriteTimeout() time.Duration {
	if GetAppConf().Server.WriteTimeout > 0 {
		return time.Duration(GetAppConf().Server.WriteTimeout) * time.Millisecond
	}
	return 2 * time.Second
}

// GetDatabaseConf 获取数据库配置
func (ac *AppConfig) GetDatabaseConf() map[string]DatabaseConf {
	if len(GetAppConf().Database) == 0 {
		return make(map[string]DatabaseConf)
	}
	return GetAppConf().Database
}

// GetDatabaseTableConf 获取数据库表配置
func (ac *AppConfig) GetDatabaseTableConf(dbName string, tableName string) TableConf {
	if len(GetAppConf().Database) == 0 {
		return TableConf{}
	}
	dbconf := GetAppConf().Database[dbName]
	tabelsConf := dbconf.TableConfig
	if len(tabelsConf) == 0 {
		return TableConf{}
	}
	for _, tableConf := range tabelsConf {
		if tableConf.TableName == tableName {
			return tableConf
		}
	}
	return TableConf{}
}

// GetLoggerConf 获取日志配置
func (ac *AppConfig) GetLoggerConf() map[string]klog.Config {
	return ac.Logger
}

// GetAuthConf 获取登录认证配置
func (ac *AppConfig) GetAuthConf() AuthConf {
	if ac.Auth.JwtSecret == "" {
		ac.Auth.JwtSecret = "cine_stream_login_jwt"
	}
	// 默认三个小时过期时间
	if ac.Auth.ExpireHours == 0 {
		ac.Auth.ExpireHours = 3
	}
	return ac.Auth
}

// GetPassportConf 获取 Passport 配置
func (ac *AppConfig) GetPassportConf() PassportConf {
	return ac.Auth.Passport
}

// GetCDNConf 获取 CDN 配置
func (ac *AppConfig) GetCDNConf() map[string]CDNConf {
	return ac.CDN
}
