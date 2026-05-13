package core

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"wails-router-link/service/types"

	"github.com/BurntSushi/toml"
)

func NewDefaultConfig() *types.AppConfig {
	pathWd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &types.AppConfig{
		AppName:   "gin-template",
		Debug:     false,
		Version:   "1.0.0",
		StaticDir: filepath.Join(pathWd, "static"),
		Host:      "localhost",
		Port:      3344,
		Sqlite: types.SqliteConfig{
			Enable:              false,
			BasePath:            filepath.Join(pathWd, "static", "db"),
			SqliteEncryptionKey: "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99",
		},
		WebSocket: types.WebSocketConfig{
			Enable:  false,
			BaseUrl: "/ws/:id",
		},
		Swagger: types.SwaggerConfig{
			Enable: false,
		},
		Cron: types.CronConfig{
			Enable: false,
		},
		Logger: types.LoggerConfig{
			Level:        "debug",
			MaxAge:       30 * 24 * time.Hour, // 保留30天
			Filepath:     filepath.Join(pathWd, "static", "logs", "app"),
			RotationTime: 24 * time.Hour, // 每天切割
		},
		Upload: types.UploadConfig{
			Filepath: filepath.Join(pathWd, "static", "uploads"),
		},
	}
}

func LoadConfig(configFile string) (*types.AppConfig, error) {
	var config *types.AppConfig
	_, err := os.Stat(configFile)
	if err != nil {
		config = NewDefaultConfig()
		config.ConfigPath = configFile
		err := SaveConfig(config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	_, err = toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, err
	}

	return config, err
}

func SaveConfig(config *types.AppConfig) error {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	if err := encoder.Encode(&config); err != nil {
		return err
	}
	return os.WriteFile(config.ConfigPath, buf.Bytes(), 0644)
}
