package sql_driver

import (
	"fmt"
	"os"
	"path/filepath"

	sqliteEncrypt "github.com/ShaoQ1ang/gorm-sqlite-cipher"
	"gorm.io/gorm"
	"wails-router-link/service/types"
	"wails-router-link/service/utility"
)

func NewSqliteDriver(config *gorm.Config, appConfig *types.AppConfig) (*gorm.DB, error) {
	// _pragma_cipher_compatibility=3
	params := fmt.Sprintf("?_pragma_key=%s&_pragma_cipher_page_size=4096", appConfig.Sqlite.SqliteEncryptionKey)
	dbFileName := utility.Tern(appConfig.Debug, "dev.db", "prod.db"+params)
	if err := os.MkdirAll(appConfig.Sqlite.BasePath, 0777); err != nil {
		panic("NewSqliteDriver: " + err.Error())
	}
	dbFullPath := filepath.Join(appConfig.Sqlite.BasePath, appConfig.AppName+"-"+dbFileName.(string))
	//db, err := gorm.Open(sqlite.Open(dbFullPath), config)
	db, err := gorm.Open(sqliteEncrypt.Open(dbFullPath), config)
	if err != nil {
		return nil, err
	}
	appConfig.Sqlite.Enable = true
	return db, nil
}

type User struct {
	ID   uint
	Name string
}

func test(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.Create(&User{Name: "Alice"})
	var user User
	db.First(&user)
	println("查询结果:", user.Name)
}
