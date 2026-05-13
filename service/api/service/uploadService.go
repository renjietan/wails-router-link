package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"wails-router-link/service/types"
)

type UploadService struct {
	db     *gorm.DB
	lock   sync.Mutex
	config *types.AppConfig
}

func NewUploadService(db *gorm.DB, appConfig *types.AppConfig) *UploadService {
	return &UploadService{
		db:     db,
		lock:   sync.Mutex{},
		config: appConfig,
	}
}

func (service *UploadService) Upload(c *gin.Context, file *multipart.FileHeader) error {
	service.lock.Lock()
	defer service.lock.Unlock()
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), "upload", ext)
	savePath := filepath.Join(service.config.Upload.Filepath, newFileName)
	return c.SaveUploadedFile(file, savePath)
}
