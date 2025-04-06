package main

import (
	"fast-gin/core"
	"fast-gin/flags"
	"fast-gin/global"
	"fmt"
)

func main() {
	flags.Run()

	var err error
	global.Config, err = core.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(global.Config.DB)

	global.Config.DB.Port = 3307
	err = core.DumpConfig()
	if err != nil {
		panic(err)
	}
}
