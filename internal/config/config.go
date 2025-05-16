package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	Address     string `yaml:"address"`
	TLSCertFile string `yaml:"tls_cert_file"`
	TLSKeyFile  string `yaml:"tls_key_file"`
}

type Config struct {
	DatabaseURL string       `yaml:"database_url"`
	Server      ServerConfig `yaml:"server"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
