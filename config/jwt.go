package config

type JWT struct {
	Expire    int    `yaml:"expire"`
	Issuer    string `yaml:"issuer"`
	SecretKey string `yaml:"secret_key"`
}
