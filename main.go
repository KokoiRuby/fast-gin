package main

import (
	"fast-gin/core"
	"fmt"
)

func main() {
	cfg, err := core.LoadConfig("./config/settings-dev.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.DB)
}
