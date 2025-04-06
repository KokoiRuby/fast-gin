package global

import (
	"fast-gin/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const VERSION = "0.0.1"

var (
	Config *config.Config
	DB     *gorm.DB
	Redis  *redis.Client
)
