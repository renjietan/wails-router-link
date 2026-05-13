package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	dto "wails-router-link/service/api/DTO"
	"wails-router-link/service/api/entity"
	"wails-router-link/service/api/service"
	"wails-router-link/service/core"
	"wails-router-link/service/types"
	"wails-router-link/service/utility/response"
)

type ConfigController struct {
	BaseController
	service        *service.ConfigService
	upload_service *service.UploadService
}

func NewConfigController(
	app *core.AppServer,
	db *gorm.DB,
	service *service.ConfigService,
	upload_service *service.UploadService,
) *ConfigController {
	return &ConfigController{
		BaseController: BaseController{
			App: app,
			DB:  db,
		},
		service:        service,
		upload_service: upload_service,
	}
}

func (config *ConfigController) RegisterConfigRouters() {
	group := config.App.Engine.Group("/config")
	group.GET("list", config.list)
	group.POST("insert", config.insert)
	group.PATCH("insertMany", config.insertMany)
	group.POST("upload", config.upload)
	group.PUT("update/:id", config.update)
	group.PUT("updates", config.updates)
	group.DELETE("delete/:id", config.delete)
	group.GET("test", config.test)
}

// @Summary 创建 配置 - 含关联 插入
// @Description 配置 描述
// @Tags Config
// @Accept json
// @Produce json
// @Param request body dto.ConfigDTO true "参数"
// @Success 200 {object} entity.ConfigEntity
// @Router /config/insert [post]
func (config *ConfigController) insert(c *gin.Context) {
	var d dto.ConfigDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.ERROR(c, err.Error())
		return
	}
	res := config.service.Insert(d)
	response.SUCCESS(c, res)
}

// @Summary 批量创建 配置
// @Description 批量配置 描述
// @Tags Config
// @Accept json
// @Produce json
// @Param request body dto.ConfigsDTO true "参数"
// @Success 200 {object} []entity.ConfigEntity
// @Router /config/insertMany [patch]
func (config *ConfigController) insertMany(c *gin.Context) {
	var d dto.ConfigsDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.ERROR(c, types.InvalidArgs)
		return
	}
	res := config.service.InsertMany(d)
	response.SUCCESS(c, res)
}

// @Summary      上传单个文件
// @Description  通过 multipart/form-data 上传文件，支持图片、文档等类型
// @Tags         Config
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "要上传的文件"
// @Success      200   {object}  types.Response  "上传成功"
// @Failure      400   {object}  types.Response  "请求参数错误"
// @Failure      500   {object}  types.Response  "服务器内部错误"
// @Router       /config/upload [post]
func (config *ConfigController) upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ERROR(c, types.InvalidArgs)
	}
	err2 := config.upload_service.Upload(c, file)
	if err2 != nil {
		response.ERROR(c, types.UploadFaild)
		return
	}
	response.SUCCESS(c, types.SuccessMsg)
}

// @Summary      更新
// @Description  根据ID更新其基本信息
// @Tags         Config
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "ID"
// @Param        req  body      dto.ConfigDTO    true  "需要更新的信息"
// @Success      200  {object}  types.Response   "更新成功"
// @Failure      400  {object}  types.Response  "请求参数错误"
// @Failure      404  {object}  types.Response  "用户不存在"
// @Failure      500  {object}  types.Response  "服务器内部错误"
// @Router       /config/update/{id} [put]
func (config *ConfigController) update(c *gin.Context) {
	id := c.Param("id")
	var d dto.ConfigDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.ERROR(c, types.InvalidArgs)
		return
	}
	if err := config.service.Update(id, d); err != nil {
		response.ERROR(c, err.Error())
		return
	}
	response.SUCCESS(c)
}

// @Summary      批量更新
// @Description  批量更新其基本信息
// @Tags         Config
// @Accept       json
// @Produce      json
// @Param        req  body      dto.UpdatesConfigDTO    true  "需要更新的信息"
// @Success      200  {object}  types.Response   "更新成功"
// @Failure      400  {object}  types.Response  "请求参数错误"
// @Failure      404  {object}  types.Response  "记录不存在"
// @Failure      500  {object}  types.Response  "服务器内部错误"
// @Router       /config/updates [put]
func (config *ConfigController) updates(c *gin.Context) {
	var d dto.UpdatesConfigDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.ERROR(c, types.InvalidArgs)
		return
	}
	if err := config.service.Updates(d.Items); err != nil {
		response.ERROR(c, err.Error())
		return
	}
	response.SUCCESS(c)
}

// @Summary      删除
// @Description  根据ID删除
// @Tags         Config
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "ID"
// @Success      200  {object}  types.Response   "更新成功"
// @Failure      400  {object}  types.Response  "请求参数错误"
// @Failure      404  {object}  types.Response  "用户不存在"
// @Failure      500  {object}  types.Response  "服务器内部错误"
// @Router       /config/delete/{id} [delete]
func (config *ConfigController) delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := config.service.Delete(id); err != nil {
		response.ERROR(c, err.Error())
		return
	}
	response.SUCCESS(c, types.SuccessMsg)
}

// @Summary      获取配置列表（分页+排序）
// @Description  分页查询配置项，按指定字段排序。默认按创建时间倒序排列。
// @Tags         Config
// @Accept       json
// @Produce      json
// @Param		 req 		query     dto.ConfigListDTO    true  "需要更新的信息"
// @Success      200        {object}  map[string]interface{} "分页数据"
// @Failure      400        {object}  map[string]interface{}  "请求参数错误"
// @Failure      500        {object}  map[string]interface{}  "服务器内部错误"
// @Router       /config/list [get]
func (config *ConfigController) list(c *gin.Context) {
	var d dto.ConfigListDTO
	if err := c.ShouldBindQuery(&d); err != nil {
		response.ERROR(c, types.InvalidArgs)
		return
	}
	list, err := config.service.List(d)
	if err != nil {
		response.ERROR(c, err.Error())
		return
	}
	response.SUCCESS(c, &list)
}

// @Summary      测试
// @Description  测试
// @Tags         Config
// @Accept       json
// @Produce      json
// @Param		 req 		query     dto.ConfigListDTO    true  "需要更新的信息"
// @Success      200        {object}  map[string]interface{} "分页数据"
// @Failure      400        {object}  map[string]interface{}  "请求参数错误"
// @Failure      500        {object}  map[string]interface{}  "服务器内部错误"
// @Router       /config/test [get]
func (config *ConfigController) test(c *gin.Context) {
	var configs []entity.ConfigEntity
	if err := config.DB.Preload("Details").Find(&configs).Error; err != nil {
		response.ERROR(c, err.Error())
		return
	}
	response.SUCCESS(c, configs)
}
