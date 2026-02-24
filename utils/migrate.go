// Package utils 数据库迁移功能
package utils

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	sqlmigrate "github.com/rubenv/sql-migrate"
)

// DatabaseConf 数据库配置（用于迁移功能，避免循环依赖）
type DatabaseConf struct {
	Host string // host
	Port int    // 端口
	Name string // 数据库名
	User string // 用户名
	Pass string // 密码
}

// RunMigrations 执行数据库迁移
// dbConfig: 数据库配置映射
// dbName: 数据库配置名称，如果为空则对所有数据库执行迁移
func RunMigrations(dbConfig map[string]DatabaseConf, dbName string) error {
	if len(dbConfig) == 0 {
		return fmt.Errorf("数据库配置为空")
	}

	// 获取迁移文件目录
	migrationsDir := filepath.Join(".", "migrations")

	// 如果指定了数据库名称，只对该数据库执行迁移
	if dbName != "" {
		dbConf, exists := dbConfig[dbName]
		if !exists {
			return fmt.Errorf("数据库配置 %s 不存在", dbName)
		}
		return runMigrationForDB(dbName, &dbConf, migrationsDir)
	}

	// 对所有数据库执行迁移
	for name, dbConf := range dbConfig {
		if err := runMigrationForDB(name, &dbConf, migrationsDir); err != nil {
			log.Printf("[RunMigrations] 数据库 %s 迁移失败: %v", name, err)
			return err
		}
	}

	return nil
}

// runMigrationForDB 对指定数据库执行迁移
func runMigrationForDB(dbName string, dbConf *DatabaseConf, migrationsDir string) error {
	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConf.User,
		dbConf.Pass,
		dbConf.Host,
		dbConf.Port,
		dbConf.Name,
	)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 设置迁移源
	migrations := &sqlmigrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	// 执行迁移
	log.Printf("[RunMigrations] 开始执行数据库 %s 的迁移...", dbName)
	n, err := sqlmigrate.Exec(db, "mysql", migrations, sqlmigrate.Up)
	if err != nil {
		return fmt.Errorf("执行迁移失败: %v", err)
	}

	if n > 0 {
		log.Printf("[RunMigrations] 数据库 %s 成功执行了 %d 个迁移", dbName, n)
	} else {
		log.Printf("[RunMigrations] 数据库 %s 没有需要执行的迁移", dbName)
	}

	return nil
}

// GetMigrationStatus 获取迁移状态
// dbConfig: 数据库配置映射
// dbName: 数据库配置名称
func GetMigrationStatus(dbConfig map[string]DatabaseConf, dbName string) ([]*sqlmigrate.MigrationRecord, error) {
	if len(dbConfig) == 0 {
		return nil, fmt.Errorf("数据库配置为空")
	}

	dbConf, exists := dbConfig[dbName]
	if !exists {
		return nil, fmt.Errorf("数据库配置 %s 不存在", dbName)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConf.User,
		dbConf.Pass,
		dbConf.Host,
		dbConf.Port,
		dbConf.Name,
	)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 获取迁移记录
	records, err := sqlmigrate.GetMigrationRecords(db, "mysql")
	if err != nil {
		return nil, fmt.Errorf("获取迁移记录失败: %v", err)
	}

	return records, nil
}
