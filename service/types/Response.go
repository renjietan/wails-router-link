package types

type BizCode int

type Response struct {
	Code    BizCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Detail  string      `json:"detail"`
}

type PaginationResponse struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

const (
	Success       = BizCode(20000)
	Failed        = BizCode(50000)
	NotAuthorized = BizCode(401)
)

const (
	SuccessMsg  = "操作成功"
	FailedMsg   = "操作失败"
	ErrorMsg    = "系统开小差了"
	InvalidArgs = "非法参数或参数解析失败"
	UploadFaild = "文件上传失败"
	NoAuth      = "无权限，清联系管理员"
)
