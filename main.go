package main

import (
	"fast-gin/core"
	"fast-gin/flags"
	"fast-gin/global"
	"fast-gin/routers"
	"fmt"
)

func main() {
	// Logging
	core.InitLogger()

	// Flags
	flags.Parse()

	// Configuration
	var err error
	global.Config, err = core.LoadConfig()
	if err != nil {
		panic(err)
	}

	// GORM
	global.DB = core.InitGorm()

	// Redis
	global.Redis = core.InitRedis()

	// Handle DB migration and version print
	flags.Run()

	// Cron (goroutine)
	//svc_cron.CronInit()

	// Gin
	routers.Run()

	fmt.Println("End of Main")
}
