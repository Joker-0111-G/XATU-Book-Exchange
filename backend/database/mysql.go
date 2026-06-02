package database

import (
	"fmt"
	"log"

	"xatu-book-exchange/config"
	"xatu-book-exchange/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL() error {
	cfg := config.AppConfig.MySQL
	dsn := cfg.DSN()

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// 自动迁移表结构
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	// 自动填充初始数据（管理员、分类、轮播图）
	SeedData()

	log.Println("MySQL 连接成功，表迁移完成")
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Book{},
		&model.Favorite{},
		&model.Order{},
		&model.Message{},
		&model.Banner{},
	)
}
