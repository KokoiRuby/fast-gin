package config

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBMode string

const (
	MYSQL  DBMode = "mysql"
	PG     DBMode = "postgres"
	SQLITE DBMode = "sqlite"
)

type DB struct {
	Mode     DBMode `yaml:"mode"` // Supports: mysql pgsql sqlite
	DBName   string `yaml:"db_name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (db DB) GetDSN() gorm.Dialector {
	switch db.Mode {
	case MYSQL:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.DBName,
		)
		return mysql.Open(dsn)
	case PG:
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.DBName,
		)
		return postgres.Open(dsn)
	case SQLITE:
		return sqlite.Open(db.DBName)
	case "":
		logrus.Warnf("Database mode not specified")
		return nil
	default:
		logrus.Fatalf("Database is not supported")
		return nil
	}
}
