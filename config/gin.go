package config

import "fmt"

type Gin struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

func (g Gin) Addr() string {
	return fmt.Sprintf("%s:%s", g.IP, g.Port)
}
