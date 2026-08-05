package config

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
