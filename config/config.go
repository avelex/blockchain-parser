package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port           int           `yaml:"port"`
	RPC            string        `yaml:"rpc"`
	BlocksInterval time.Duration `yaml:"blocks_interval"`
	StartBlock     int           `yaml:"start_block,omitempty"`
}

func LoadConfig(path string) (Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
