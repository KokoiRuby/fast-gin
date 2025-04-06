package core

import (
	"fast-gin/config"
	"fast-gin/flags"
	"fast-gin/global"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() (cfg *config.Config, err error) {
	cfg = new(config.Config)
	file, err := os.Open(flags.Options.File)
	if err != nil {
		return cfg, fmt.Errorf("error when reading configuration file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("error when decoding YAML: %w", err)
	}
	return
}

func DumpConfig() error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		return fmt.Errorf("error when dumping configuration: %w", err)
	}
	err = os.WriteFile(flags.Options.File, byteData, 0666)
	if err != nil {
		return fmt.Errorf("error when dumping configuration: %w", err)
	}
	fmt.Println("Configuration dumped successfully")
	return nil
}
