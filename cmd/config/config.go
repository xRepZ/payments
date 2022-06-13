package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Http *ServerConfig `yaml:"http"`
	} `yaml:"server"`

	Storage struct {
		Postgres *DbConfig `yaml:"postgres"`
	} `yaml:"storage"`
}

type ServerConfig struct {
	Listen string `yaml:"listen"`
}

type DbConfig struct {
	Dsn string `yaml:"dsn"`
}

func ParseConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open config: %w", err)
	}

	c := &Config{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal config: %w", err)
	}

	return c, nil
}
