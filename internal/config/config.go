package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
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

func Load(configFileName string, address string, idleTimeout time.Duration, maxConnections int, MaxMessageSize string) (*Config, error) {
	if configFileName == "" {
		return nil, errors.New("empty config file name")
	}
	dataCfg, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatal(err)
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

	if config.Network.Address == "" {
		config.Network.Address = address
	}
	if config.Network.MaxConnections == 0 {
		config.Network.MaxConnections = maxConnections
	}
	if config.Network.MaxMessageSize == "" {
		config.Network.MaxMessageSize = MaxMessageSize
	}
	if config.Network.IdleTimeout == 0 {
		config.Network.IdleTimeout = idleTimeout
	}

	return &config, nil
}
