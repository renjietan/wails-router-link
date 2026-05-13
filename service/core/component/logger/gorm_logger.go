package logger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger 实现gorm的logger接口，将SQL输出到logrus
type GormLogger struct {
	*logrus.Logger
	SlowThreshold time.Duration // 慢查询阈值
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithContext(ctx).Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithContext(ctx).Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WithContext(ctx).Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := logrus.Fields{
		"rows":     rows,
		"duration": elapsed,
		"SQL":      sql,
	}
	if err != nil {
		fields["error"] = err
		l.Logger.WithContext(ctx).WithFields(fields).Error("❌ GOORM")
		return
	}
	if elapsed > l.SlowThreshold && l.SlowThreshold != 0 {
		l.Logger.WithContext(ctx).WithFields(fields).Warn("⚠️ GOORM")
	} else {
		l.Logger.WithContext(ctx).WithFields(fields).Debug("🔵 GOORM")
	}
}
