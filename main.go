package main

import (
	"fast-gin/core"
	"fast-gin/flags"
	"fast-gin/global"
	"github.com/sirupsen/logrus"
)

func main() {
	core.InitLogger()

	flags.Run()

	var err error
	global.Config, err = core.LoadConfig()
	if err != nil {
		panic(err)
	}
	logrus.Infof("Configuration loaded successfully: %v", global.Config.DB)
}
