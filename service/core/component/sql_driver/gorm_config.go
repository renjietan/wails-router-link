package sql_driver

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"wails-router-link/service/core/component/logger"
)

// NewGormConfig Mortal 2026-05-02 00:21:24 CST 将 GORM 配置 单独拎出来 方便其他类型数据库共用
func NewGormConfig(log *logrus.Logger) *gorm.Config {
	return &gorm.Config{
		Logger: NewGormLogger(log),
		//Logger: logger2.GetLogger().Desugar(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "gin_", // 设置表前缀
			SingularTable: false,  // 使用单数表名形式
		},
	}
}

func NewGormLogger(log *logrus.Logger) *logger.GormLogger {
	gormLogger := &logger.GormLogger{
		Logger:        log,
		SlowThreshold: 200 * time.Millisecond,
	}
	return gormLogger
}
