package config

import (
	"os"
	"strconv"

	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config Global configuration structure
type Config struct {
	Server *ServerConfig `yaml:"server"`
	Logger *LoggerConfig `yaml:"logger"`
	API    *APIConfig    `yaml:"api"`
}

// APIConfig API configuration
type APIConfig struct {
	// API base URL
	BaseURL string `yaml:"base_url"`
	// API token
	Token string `yaml:"token"`
	// API timeout
	Timeout int `yaml:"timeout"`
	// API debug mode
	Debug bool `yaml:"debug"`
	// API insecure skip verify
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
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

var config *Config

// getEnvString Get string value from environment variable, return default value if not exists
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt Get integer value from environment variable, return default value if not exists or cannot be parsed
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Load Load configuration file
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "read config file failed")
	}

	config = &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return errors.Wrap(err, "unmarshal config failed")
	}

	// Override configuration from environment variables
	if config.Server != nil {
		config.Server.Host = getEnvString("LISTEN_ADDRESS", config.Server.Host)
		config.Server.Port = getEnvInt("LISTEN_PORT", config.Server.Port)
	}

	if config.API != nil {
		config.API.BaseURL = getEnvString("SAFELINE_ADDRESS", config.API.BaseURL)
		config.API.Token = getEnvString("SAFELINE_API_TOKEN", config.API.Token)
	}

	return nil
}

// GetServer Get server configuration
func GetServer() *ServerConfig {
	if config == nil {
		return nil
	}
	return config.Server
}

// GetLogger Get logger configuration
func GetLogger() *LoggerConfig {
	if config == nil {
		return nil
	}
	return config.Logger
}

// GetAPI Get API configuration
func GetAPI() *APIConfig {
	if config == nil {
		return nil
	}
	return config.API
}
