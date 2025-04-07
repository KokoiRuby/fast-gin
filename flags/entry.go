package flags

import (
	"fast-gin/global"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type FlagOptions struct {
	File      string
	Version   bool
	DB        bool
	Resource  string // user
	Operation string // create, list, remove
}

var Options FlagOptions

func Parse() {
	flag.StringVar(&Options.File, "f", "./config/settings.yaml", "Configuration file")
	flag.StringVar(&Options.Resource, "res", "", "Resource: [user]")
	flag.StringVar(&Options.Operation, "op", "", "Operation: [create|list|remove]")
	flag.BoolVar(&Options.Version, "v", false, "Print version information")
	flag.BoolVar(&Options.DB, "db", false, "Database migration")
	flag.Parse()
}

func Run() {
	if Options.DB {
		MigrateDB()
		os.Exit(0)
	}
	if Options.Version {
		fmt.Println("Version:", global.VERSION)
		os.Exit(0)
	}
	if Options.Resource == "user" {
		var user User
		switch Options.Operation {
		case "create":
			user.Create()
		case "list":
			user.List()
		case "remove":
			user.Remove()
		default:
			logrus.Fatalf("Operation [%s] not supported", Options.Operation)
		}
		os.Exit(0)
	}
}
