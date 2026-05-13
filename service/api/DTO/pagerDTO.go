package dto

type PagerDTO struct {
	// 当前页码 默认: 1
	Page int `form:"page" json:"page" example:"1" binding:"omitempty,min=1"`
	// 条数 默认: 10
	PageSize int `form:"page_size" json:"page_size" example:"10" binding:"omitempty,min=1,max=100"`
	// 排序 默认: created_at
	SortBy string `form:"sort_by" json:"sort_by" example:"created_at"\`
	// 排序/规则 默认:desc（倒序）
	Order string `form:"order" json:"order" binding:"omitempty,oneof=asc desc"`
}
