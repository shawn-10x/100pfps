package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func Connect() {
	var err error

	logger_ := logger.Info
	if os.Getenv("MODE") == "Release" {
		logger_ = logger.Silent
	}
	db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger_),
	})
	if err != nil {
		panic("failed to connect database")
	}
}
