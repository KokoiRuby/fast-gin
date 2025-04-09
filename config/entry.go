package config

type Config struct {
	DB     DB     `yaml:"db"`
	Redis  Redis  `yaml:"redis"`
	Gin    Gin    `yaml:"gin"`
	JWT    JWT    `yaml:"jwt"`
	Upload Upload `yaml:"upload"`
}
