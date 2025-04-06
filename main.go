package main

import (
	"fast-gin/core"
	"fast-gin/flags"
	"fast-gin/global"
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

	fmt.Println("End of Main")
}
