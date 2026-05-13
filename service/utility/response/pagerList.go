package response

import (
	"strings"

	"gorm.io/gorm"
	dto "wails-router-link/service/api/DTO"
	"wails-router-link/service/types"
)

// Paginate 参数:
//   - db: GORM 查询链，已经包含了模型和 Where 条件
//   - req: 分页请求参数
//   - allowedSortFields: 允许排序的字段白名单，key 为字段名（如 "created_at"）
//   - resultSlice: 用于接收查询结果的切片指针，例如 &[]entity.xxxx{}

func Paginate(db *gorm.DB, req dto.PagerDTO, allowedSortFields map[string]bool, res interface{}) (*types.PaginationResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	if allowedSortFields == nil {
		allowedSortFields = map[string]bool{
			"id":         true,
			"created_at": true,
			"updated_at": true,
		}
	}
	if !allowedSortFields[req.SortBy] {
		// 排序字段非法，返回错误
		return nil, gorm.ErrInvalidField
	}

	orderDirection := strings.ToUpper(req.Order)
	var total int64
	// 计算总数（会自动过滤软删除）
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 执行分页查询
	offset := (req.Page - 1) * req.PageSize
	err := db.
		Order(req.SortBy + " " + orderDirection).
		Limit(req.PageSize).
		Offset(offset).
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	// 构造响应
	return &types.PaginationResponse{
		List:       res,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
