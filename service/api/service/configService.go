package service

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
	dto "wails-router-link/service/api/DTO"
	"wails-router-link/service/api/entity"
	"wails-router-link/service/types"
	"wails-router-link/service/utility/response"
)

type ConfigService struct {
	db   *gorm.DB
	lock sync.Mutex
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db, lock: sync.Mutex{}}
}

func (c *ConfigService) Insert(d dto.ConfigDTO) entity.ConfigEntity {
	c.lock.Lock()
	defer c.lock.Unlock()
	res := entity.ConfigEntity{
		Name:  d.Name,
		Value: d.Value,
		Details: []entity.ConfigDetailEntity{
			{Content: d.Details.Content},
		},
	}
	c.db.Create(&res)
	return res
}

func (c *ConfigService) InsertMany(d dto.ConfigsDTO) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	var res []entity.ConfigEntity
	for _, v := range d.Items {
		temp := entity.ConfigEntity{
			Name:  v.Name,
			Value: v.Value,
		}
		res = append(res, temp)
	}
	tx := c.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败: %w", tx.Error)
	}
	if err := tx.CreateInBatches(res, 100); err.Error != nil {
		tx.Rollback()
		return err.Error
	}
	return tx.Commit().Error
}

func (c *ConfigService) Update(id string, configDTO dto.ConfigDTO) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	var config entity.ConfigEntity
	if err := c.db.First(&config, id).Error; err != nil {
		return err
	}
	config.Name = configDTO.Name
	config.Value = configDTO.Value
	c.db.Save(&config).Scan(&config)
	return nil
}

func (c *ConfigService) Updates(dtos []entity.ConfigEntity) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	var ids []uint
	for _, u := range dtos {
		ids = append(ids, u.ID)
	}
	// 构建后的SQL类似：UPDATE users SET name = CASE id WHEN 1 THEN 'name1' WHEN 2 THEN 'name2' END WHERE id IN (1,2)
	var nameCases, valueCases string
	for _, u := range dtos {
		nameCases += fmt.Sprintf("WHEN %d THEN '%s' ", u.ID, u.Name)
		valueCases += fmt.Sprintf("WHEN %d THEN '%s' ", u.ID, u.Value)
	}
	sql := fmt.Sprintf(`
		UPDATE config SET 
            name = CASE id %s END,
            value = CASE id %s END,
            updated_at = NOW()
        WHERE id IN (?)
    `, nameCases, valueCases)
	return c.db.Exec(sql, ids).Error
}

func (c *ConfigService) Delete(id string) (entity.ConfigEntity, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var a entity.ConfigEntity
	c.db.Delete(&a, id)
	return a, nil
}

func (c *ConfigService) List(dto dto.ConfigListDTO) (*types.PaginationResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	query := c.db.Preload("Details").Where("name LIKE ?", "%"+dto.Name+"%").Preload("Details").Find(&entity.ConfigEntity{})
	var configs []entity.ConfigEntity
	res, err := response.Paginate(query, dto.PagerDTO, nil, configs)
	if err != nil {
		return nil, err
	}
	return res, nil
}
