package main

import (
	"fast-gin/core"
	"fast-gin/flags"
	"fast-gin/global"
	"github.com/sirupsen/logrus"
)

func main() {
	// Logging
	core.InitLogger()

	// Flags
	flags.Run()

	// Configuration
	var err error
	global.Config, err = core.LoadConfig()
	if err != nil {
		panic(err)
	}
	logrus.Infof("DB Configuration: %v", global.Config.DB)
	logrus.Errorf("Test Configuration error")

	// GORM
	global.DB = core.InitGorm()

	// Redis
	global.Redis = core.InitRedis()

}
