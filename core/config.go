package core

import (
	"fast-gin/config"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(filename string) (cfg *config.Config, err error) {
	cfg = new(config.Config)

	file, err := os.Open(filename)
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
