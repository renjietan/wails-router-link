package types

import "time"

type AppConfig struct {
	AppName    string `toml:"app_name"`
	Debug      bool   `toml:"-"`
	Version    string `toml:"version"`
	ConfigPath string `toml:"-"`
	StaticDir  string `toml:"static_dir"`
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	LogDir     string `toml:"log_dir"`
	/* 注意这里不可以用匿名嵌入子结构体，否则写入toml文件时, 无法形成树状结构 */
	Sqlite    SqliteConfig    `toml:"sqlite"`
	Logger    LoggerConfig    `toml:"logger"`
	Upload    UploadConfig    `toml:"upload"`
	WebSocket WebSocketConfig `toml:"websocket"`
	Cron      CronConfig      `toml:"cron"`
	Swagger   SwaggerConfig   `toml:"swagger"`
}
type MysqlConfig struct {
	Enable   bool   `toml:"enable"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"data_base"`
}

type SqliteConfig struct {
	Enable              bool   `toml:"enable"`                // 是否启用
	BasePath            string `toml:"base_path"`             // db文件基础路径
	SqliteEncryptionKey string `toml:"sqlite_encryption_key"` // 加密密钥
}

type LoggerConfig struct {
	Level        string        `toml:"level"`
	Filepath     string        `toml:"file_path"`
	MaxAge       time.Duration `toml:"max_age"`
	RotationTime time.Duration `toml:"rotation_time"`
}

type UploadConfig struct {
	Filepath string `toml:"file_path"`
}

type WebSocketConfig struct {
	Enable  bool   `toml:"enable"`
	BaseUrl string `toml:"base_url"`
}

type CronConfig struct {
	Enable bool `toml:"enable"`
}

type SwaggerConfig struct {
	Enable bool `toml:"enable"`
}
