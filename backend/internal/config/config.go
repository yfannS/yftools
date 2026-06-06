package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"md2html/pkg/logger"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Mode    string `mapstructure:"mode"` // debug / release
	LogLevel string `mapstructure:"log_level"` // DEBUG / INFO / WARN / ERROR
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"` // e.g. "168h" = 7 days
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 环境变量覆盖（可选，优先级高于 yaml）
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Fatal("Failed to unmarshal config: %v", err)
	}

	return &cfg
}

func (c *Config) InitDB() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
	logger.Info("Connecting to database: %s@tcp(%s:%d)/%s", c.Database.User, c.Database.Host, c.Database.Port, c.Database.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("Failed to open database: %v", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database: %v", err)
	}

	logger.Info("Database connected successfully")
	return db
}
