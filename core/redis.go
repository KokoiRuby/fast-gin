package core

import (
	"context"
	"fast-gin/global"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func InitRedis() *redis.Client {
	cfg := global.Config
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Errorf("Failed to connect to redis: %s", err)
		return nil
	}
	logrus.Infof("Connect to redis successfully")
	return rdb
}
