package flags

import (
	"flag"
)

type FlagOptions struct {
	File string
}

var Options FlagOptions

func Run() {
	flag.StringVar(&Options.File, "f", "./config/settings.yaml", "configuration file")
	flag.Parse()
}
