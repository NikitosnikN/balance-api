package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	exampleConfigPath = "config.example.yaml"
)

type ServerConfig struct {
	Port uint `yaml:"port"`
}

type RpcConfig struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type Config struct {
	HTTPPort       uint          `yaml:"port"`
	MetricsPort    uint          `yaml:"metricsPort"`
	WorkerInterval time.Duration `yaml:"workerInterval"`
	Rpcs           []RpcConfig
}

func openExampleConfig() (*os.File, error) {
	log.Printf("Loading default config (path %s)\n", exampleConfigPath)

	f, err := os.Open(filepath.Clean(exampleConfigPath))

	if err != nil {
		return nil, fmt.Errorf("failed to open example config file: %w", err)
	}

	return f, nil

}

func LoadConfigFromFile(path string) (*Config, error) {
	var cfg *Config

	f, err := os.Open(filepath.Clean(path))

	if err != nil && os.IsNotExist(err) {
		f, err = openExampleConfig()
	}

	if err != nil {
		return nil, err
	}

	err = yaml.NewDecoder(f).Decode(&cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
