package core

import (
	"fast-gin/global"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func InitGorm() (db *gorm.DB) {
	cfg := global.Config.DB

	dialector := cfg.GetDSN()
	if dialector == nil {
		return
	}

	// Open initialize db session based on dialector
	database, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Get DB connection pool
	sqlDB, err := database.DB()
	if err != nil {
		logrus.Fatalf("Failed to get database connection pool: %s", err)
		return
	}
	err = sqlDB.Ping()
	if err != nil {
		logrus.Fatalf("Failed to probe database connection pool liveness: %s", err)
		return
	}

	// Configure DB connection pool
	// TODO: Add to configuration file
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logrus.Infof("DB initialized successfully")
	return
}
