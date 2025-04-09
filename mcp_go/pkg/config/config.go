package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config Global configuration structure
type Config struct {
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
}

// ServerConfig Server configuration
type ServerConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
	Secret  string `yaml:"secret"`
}

// LoggerConfig Logger configuration
type LoggerConfig struct {
	Level       string `yaml:"level"`
	FilePath    string `yaml:"file_path"`
	Console     bool   `yaml:"console"`
	Caller      bool   `yaml:"caller"`
	Development bool   `yaml:"development"`
}

var (
	globalConfig *Config
)

// Load Load configuration from file
func Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}

	globalConfig = config
	return nil
}

// Get Get global configuration
func Get() *Config {
	return globalConfig
}

// GetServer Get server configuration
func GetServer() ServerConfig {
	if globalConfig == nil {
		return ServerConfig{}
	}
	return globalConfig.Server
}

// GetLogger Get logger configuration
func GetLogger() LoggerConfig {
	if globalConfig == nil {
		return LoggerConfig{}
	}
	return globalConfig.Logger
}
