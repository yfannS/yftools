package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"md2html/pkg/logger"
	appRedis "md2html/pkg/redis"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type ServerConfig struct {
	Port     string `mapstructure:"port"`
	Mode     string `mapstructure:"mode"`      // debug / release
	LogLevel string `mapstructure:"log_level"` // DEBUG / INFO / WARN / ERROR
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"` // e.g. "168h" = 7 days
}

type RateLimitConfig struct {
	Enabled  bool                    `mapstructure:"enabled"`
	Prefix   string                  `mapstructure:"prefix"`
	Register RegisterRateLimitConfig `mapstructure:"register"`
	Login    LoginRateLimitConfig    `mapstructure:"login"`
}

type RegisterRateLimitConfig struct {
	IPPerMinute     LimitRuleConfig `mapstructure:"ip_per_minute"`
	IPPerHour       LimitRuleConfig `mapstructure:"ip_per_hour"`
	UsernamePerHour LimitRuleConfig `mapstructure:"username_per_hour"`
}

type LoginRateLimitConfig struct {
	IPPerMinute         LimitRuleConfig `mapstructure:"ip_per_minute"`
	UsernamePer15m      LimitRuleConfig `mapstructure:"username_per_15m"`
	FailUsernamePer15m  LimitRuleConfig `mapstructure:"fail_username_per_15m"`
	FailUsernameIPPer5m LimitRuleConfig `mapstructure:"fail_username_ip_per_5m"`
}

type LimitRuleConfig struct {
	Limit  int    `mapstructure:"limit"`
	Window string `mapstructure:"window"`
}

func (c RateLimitConfig) KeyPrefix() string {
	if c.Prefix == "" {
		return "md2html:auth:ratelimit"
	}
	return c.Prefix
}

func (c LimitRuleConfig) Duration() time.Duration {
	if c.Window == "" {
		return 0
	}

	window, err := time.ParseDuration(c.Window)
	if err != nil || window <= 0 {
		return 0
	}

	return window
}

func (c LimitRuleConfig) IsEnabled() bool {
	return c.Limit > 0 && c.Duration() > 0
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&tls=false",
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

// InitRedis 初始化 Redis 连接
func (c *Config) InitRedis() {
	if err := appRedis.Init(appRedis.Config{
		Host: c.Redis.Host,
		Port: c.Redis.Port,
	}); err != nil {
		logger.Warn("Redis connection failed, session will use JWT-only: %v", err)
	}
}
