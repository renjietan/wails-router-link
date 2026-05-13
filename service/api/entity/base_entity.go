package entity

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	// 主键，类型为 bigint unsigned，自增
	ID uint `gorm:"primarykey;autoIncrement" json:"id"`
	// 时间字段，使用 datetime(3) 以支持毫秒精度
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;" json:"updated_at"`
	// 软删除字段，使用指针，允许为 NULL
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
}
