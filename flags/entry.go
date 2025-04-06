package flags

import (
	"fast-gin/global"
	"flag"
	"fmt"
)

type FlagOptions struct {
	File    string
	Version bool
	DB      bool
}

var Options FlagOptions

func Parse() {
	flag.StringVar(&Options.File, "f", "./config/settings.yaml", "Configuration file")
	flag.BoolVar(&Options.Version, "v", false, "Print version information")
	flag.BoolVar(&Options.DB, "db", false, "Database migration")
	flag.Parse()
}

func Run() (ok bool) {
	if Options.DB {
		fmt.Println("Database migration")
		return true
	}
	if Options.Version {
		fmt.Println("Version:", global.VERSION)
		return true
	}
	return false
}
