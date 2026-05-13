package service

import (
	"gorm.io/gorm"
	"wails-router-link/service/api/entity"
	"wails-router-link/service/types"
)

type MigrationService struct {
	db        *gorm.DB
	appConfig *types.AppConfig
}

func NewMigrationService(db *gorm.DB, appConfig *types.AppConfig) *MigrationService {
	return &MigrationService{
		db:        db,
		appConfig: appConfig,
	}
}

func (s *MigrationService) StartMigrate() {
	go func() {
		err := s.TableMigration()
		if err != nil {
			return
		}
	}()
}

func (s *MigrationService) TableMigration() error {
	var err error
	err = s.db.AutoMigrate(&entity.ConfigEntity{})
	if err != nil {
		return err
	}
	err = s.db.AutoMigrate(&entity.ConfigDetailEntity{})
	if err != nil {
		return err
	}
	return nil
}
