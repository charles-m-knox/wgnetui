package database

import (
	"database/sql"
	"log"
	"time"

	"wgnetui/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// DB is for higher level Go operations
	DB *gorm.DB
	// SQLDB is for more low-level SQL operations
	SQLDB          *sql.DB
	OpenedFileName string
	OpenedFilePath string
)

func Reconnect() {
	if DB == nil || SQLDB == nil {
		return
	}

	err := SQLDB.Close()
	if err != nil {
		log.Fatalf(
			"reconnect close error: %v",
			err.Error(),
		)
	}

	DB = Connect(OpenedFilePath)
}

func Connect(file string) *gorm.DB {
	var err error

	DB, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err.Error())
	}

	SQLDB = sqlDB

	// may help with db locking in larger queries
	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	DB.Exec("PRAGMA busy_timeout=60000;")
	DB.Exec("PRAGMA journal_mode=WAL;")

	err = DB.AutoMigrate(
		&models.WgConfig{},
		&models.GenerationForm{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err.Error())
	}

	return DB
}
