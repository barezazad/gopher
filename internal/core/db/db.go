package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(DSN string, debug bool) (DB *gorm.DB) {
	var err error
	var logLevel logger.LogLevel
	switch debug {
	case false:
		logLevel = logger.Silent
	case true:
		logLevel = logger.Info
	default:
		logLevel = logger.Silent
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 10 * time.Second, // Slow SQL threshold
			LogLevel:      logLevel,         // Log level
			Colorful:      true,             // Disable color
		},
	)

	DB, err = gorm.Open(mysql.Open(DSN),
		&gorm.Config{
			Logger: newLogger,
		})

	if err != nil {
		log.Fatalln(err.Error())
	}

	return
}

func ConnectActivityDB(activityDSN string) (db *gorm.DB) {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	db, err = gorm.Open(mysql.Open(activityDSN),
		&gorm.Config{
			Logger: newLogger,
		})
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}
