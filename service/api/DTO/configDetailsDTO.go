package dto

type ConfigDetailsDTO struct {
	// 配置详情
	Content string `json:"content" binding:"required" example:"配置详情"`
}
