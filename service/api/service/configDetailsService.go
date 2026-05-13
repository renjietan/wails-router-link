package service

import (
	"sync"

	"gorm.io/gorm"
	"wails-router-link/service/api/entity"
)

type ConfigDetailsService struct {
	db   *gorm.DB
	lock sync.Mutex
}

func NewConfigDetailService(db *gorm.DB) *ConfigDetailsService {
	return &ConfigDetailsService{db: db, lock: sync.Mutex{}}
}

func (c *ConfigDetailsService) GetById(id string) []entity.ConfigDetailEntity {
	var details []entity.ConfigDetailEntity
	c.db.Joins("Config").Find(&details).Where("id = ?", id)
	return details
}
