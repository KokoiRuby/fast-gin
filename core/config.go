package core

import (
	"fast-gin/config"
	"fast-gin/flags"
	"fast-gin/global"
	"github.com/sirupsen/logrus"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() (cfg *config.Config, err error) {
	cfg = new(config.Config)
	file, err := os.Open(flags.Options.File)
	if err != nil {
		logrus.Fatalf("error when reading configuration file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		logrus.Fatalf("error when decoding YAML: %s", err)
	}
	logrus.Infof("Configuration [%s] loaded successfully", flags.Options.File)
	return
}

func DumpConfig() error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		logrus.Errorf("error when dumping configuration: %s", err)
		return err
	}
	err = os.WriteFile(flags.Options.File, byteData, 0666)
	if err != nil {
		logrus.Errorf("error when dumping configuration: %s", err)
		return err
	}
	logrus.Infof("Configuration [%s] dumped successfully", flags.Options.File)
	return nil
}
