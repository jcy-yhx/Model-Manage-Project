package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用程序总配置
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	JWT       JWTConfig       `yaml:"jwt"`
	Quota     QuotaConfig     `yaml:"quota"`
	Simulator SimulatorConfig `yaml:"simulator"`
}

// ServerConfig HTTP 服务配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig MySQL 配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

// DSN 返回 MySQL 连接字符串
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	Secret          string `yaml:"secret"`
	AccessTokenTTL  string `yaml:"access_token_ttl"`
	RefreshTokenTTL string `yaml:"refresh_token_ttl"`
}

// QuotaConfig 配额配置
type QuotaConfig struct {
	DefaultDailyTokens int `yaml:"default_daily_tokens"`
}

// SimulatorConfig 模拟器配置
type SimulatorConfig struct {
	MinInterval string `yaml:"min_interval"`
	MaxInterval string `yaml:"max_interval"`
	Enabled     bool   `yaml:"enabled"`
}

// Load 加载配置：读取 YAML → 环境变量覆盖
func Load() (*Config, error) {
	path := configPath()

	cfg := defaultConfig()
	if err := loadYAML(path, cfg); err != nil {
		return nil, fmt.Errorf("load config %s: %w", path, err)
	}

	applyEnvOverrides(cfg)
	return cfg, nil
}

func configPath() string {
	path := "config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		return p
	}
	flag.StringVar(&path, "config", "config.yaml", "path to config file")
	flag.Parse()
	return path
}

func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "root",
			Password: "root123",
			DBName:   "ai_gateway",
			Charset:  "utf8mb4",
		},
		Redis: RedisConfig{
			Addr: "127.0.0.1:6379",
			DB:   0,
		},
		JWT: JWTConfig{
			Secret:          "ai-gateway-jwt-secret-change-in-production",
			AccessTokenTTL:  "24h",
			RefreshTokenTTL: "168h",
		},
		Quota: QuotaConfig{
			DefaultDailyTokens: 1000000,
		},
		Simulator: SimulatorConfig{
			MinInterval: "3s",
			MaxInterval: "6s",
			Enabled:     true,
		},
	}
}

func loadYAML(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

// applyEnvOverrides 环境变量覆盖 YAML 配置
func applyEnvOverrides(cfg *Config) {
	// 数据库连接：DB_DSN 整串覆盖（Docker Compose 方式）
	if v := os.Getenv("DB_DSN"); v != "" {
		// 解析 DSN 填充 DatabaseConfig（简单解析，仅覆盖关键字段）
		// 由于 GORM 直接使用 DSN 字符串，这里仅记录 override 标记
		// 实际使用由 DatabaseConfig.DSN() 统一处理
		cfg.Database.Host = "__dsn_override__"
	}

	if v := os.Getenv("REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Server.Port)
	}
	if v := os.Getenv("GIN_MODE"); v != "" {
		cfg.Server.Mode = v
	}
}
