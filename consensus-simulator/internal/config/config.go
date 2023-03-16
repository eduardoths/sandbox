package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	YAML_FILE = "raft.yaml"
)

type Config struct {
	Servers ServerConfig `yaml:"servers"`
}

var globalConfig *Config

type ServerConfig struct {
	Leader    string   `yaml:"leader"`
	Followers []string `yaml:"followers"`
}

func init() {
	f, err := os.ReadFile(YAML_FILE)
	if err != nil {
		defaults()
		return
	}
	globalConfig = &Config{}
	if err := yaml.Unmarshal(f, globalConfig); err != nil {
		defaults()
	}
}

func defaults() {
	globalConfig = &Config{
		Servers: ServerConfig{
			Leader:    "localhost:8080",
			Followers: []string{},
		},
	}
}

func Get() Config {
	return *globalConfig
}
