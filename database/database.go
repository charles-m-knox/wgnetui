package database

import (
	"log"

	"wgnetui/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	OpenedFileName string
	OpenedFilePath string
)

func Connect(file string) *gorm.DB {
	var err error

	DB, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err.Error())
	}

	err = DB.AutoMigrate(
		&models.WgConfig{},
		&models.GenerationForm{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err.Error())
	}

	return DB
}
