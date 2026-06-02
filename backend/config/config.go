package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Upload UploadConfig `mapstructure:"upload"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.DBName, m.Charset)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret               string `mapstructure:"secret"`
	AccessExpireStr      string `mapstructure:"access_expire"`
	RefreshExpireStr     string `mapstructure:"refresh_expire"`
	AccessExpire         time.Duration
	RefreshExpire        time.Duration
}

type UploadConfig struct {
	MaxSize    int      `mapstructure:"max_size"`
	Dir        string   `mapstructure:"dir"`
	AllowTypes []string `mapstructure:"allow_types"`
}

var AppConfig *Config

func InitConfig() error {
	workDir, _ := os.Getwd()
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(workDir + "/config")
	v.AddConfigPath(workDir)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := v.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 将字符串格式的时间解析为 time.Duration
	AppConfig.JWT.AccessExpire, _ = time.ParseDuration(AppConfig.JWT.AccessExpireStr)
	AppConfig.JWT.RefreshExpire, _ = time.ParseDuration(AppConfig.JWT.RefreshExpireStr)
	if AppConfig.JWT.AccessExpire == 0 {
		AppConfig.JWT.AccessExpire = 2 * time.Hour
	}
	if AppConfig.JWT.RefreshExpire == 0 {
		AppConfig.JWT.RefreshExpire = 7 * 24 * time.Hour
	}

	return nil
}
