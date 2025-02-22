package config

import (
	"bytes"
	"concurrency_hw1/pkg/logger"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Network *NetworkConfig `yaml:"network"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

func Load(log *logger.Logger, configFileName string) (*Config, error) {
	if configFileName == "" {
		return nil, errors.New("empty config file name")
	}
	dataCfg, err := os.ReadFile(configFileName)
	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to read config file")
	}

	reader := bytes.NewReader(dataCfg)

	if reader == nil {
		return nil, errors.New("incorrect reader")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("falied to read buffer")
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
