package dto

import "wails-router-link/service/api/entity"

type ConfigDTO struct {
	// 配置名称
	Name string `json:"name" binding:"required" example:"name-张三"`
	// 配置值
	Value   string           `json:"value" binding:"required" example:"value-张三"`
	Details ConfigDetailsDTO `json:"details" binding:"required"`
}

type ConfigsDTO struct {
	Items []ConfigDTO `json:"items"`
}

type UpdatesConfigDTO struct {
	Items []entity.ConfigEntity `json:"items"`
}

type ConfigListDTO struct {
	PagerDTO
	// 名称
	Name string `form:"name" json:"name" example:"张三" description:"名称"`
}
