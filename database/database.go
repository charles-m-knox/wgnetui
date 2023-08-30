package database

import (
	"database/sql"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"wgnetui/constants"
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

// Initialize is the function that runs upon app startup only. OpenedFilePath is
// a global variable (contained in this package) that can be set using command
// line flags - it will use that value to open the database if specified.
func Initialize() {
	if OpenedFilePath != "" {
		OpenedFileName = filepath.Base(OpenedFilePath)
		Connect(OpenedFilePath)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working dir: %v", err.Error())
	}

	OpenedFileName = constants.DefaultFileName
	OpenedFilePath = path.Join(cwd, constants.DefaultFileName)

	Connect(OpenedFilePath)
}

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

	log.Printf("db opening %v", file)

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
		&models.GenerationForm{},
		&models.WgConfig{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err.Error())
	}

	return DB
}
