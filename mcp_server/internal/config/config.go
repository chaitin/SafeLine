package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config is the complete MCP server configuration.
type Config struct {
	Server    *ServerConfig     `yaml:"server"`
	Logger    *LoggerConfig     `yaml:"logger"`
	Instances []*InstanceConfig `yaml:"instances"`
}

// InstanceConfig identifies one independently authenticated SafeLine instance.
// The API token is deliberately referenced through a file instead of being
// embedded in the main configuration file.
type InstanceConfig struct {
	ID                 string `yaml:"id"`
	DisplayName        string `yaml:"display_name"`
	BaseURL            string `yaml:"base_url"`
	TokenFile          string `yaml:"token_file"`
	Timeout            int    `yaml:"timeout"`
	Debug              bool   `yaml:"debug"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

// ServerConfig is the MCP server listener configuration.
type ServerConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
}

// LoggerConfig is the logging configuration.
type LoggerConfig struct {
	Level       string `yaml:"level"`
	FilePath    string `yaml:"file_path"`
	Console     bool   `yaml:"console"`
	Caller      bool   `yaml:"caller"`
	Development bool   `yaml:"development"`
}

var config *Config

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Load reads and validates the configuration file. Relative token_file paths
// are resolved from the directory containing the configuration file so that a
// deployment does not depend on the server process working directory.
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "read config file failed")
	}

	next := &Config{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(next); err != nil {
		return errors.Wrap(err, "unmarshal config failed")
	}

	if next.Server != nil {
		next.Server.Host = getEnvString("LISTEN_ADDRESS", next.Server.Host)
		next.Server.Port = getEnvInt("LISTEN_PORT", next.Server.Port)
	}

	if next.Server == nil {
		return errors.New("server configuration is required")
	}
	if strings.TrimSpace(next.Server.Name) == "" {
		return errors.New("server.name is required")
	}
	if strings.TrimSpace(next.Server.Version) == "" {
		return errors.New("server.version is required")
	}
	if strings.TrimSpace(next.Server.Host) == "" {
		return errors.New("server.host is required")
	}
	if next.Server.Port < 1 || next.Server.Port > 65535 {
		return errors.New("server.port must be between 1 and 65535")
	}
	if next.Logger == nil {
		return errors.New("logger configuration is required")
	}

	if err := normalizeInstances(next, filepath.Dir(path)); err != nil {
		return err
	}

	config = next
	return nil
}

func normalizeInstances(cfg *Config, configDir string) error {
	if len(cfg.Instances) == 0 {
		return errors.New("at least one SafeLine instance is required")
	}

	seen := make(map[string]struct{}, len(cfg.Instances))
	for index, instance := range cfg.Instances {
		if instance == nil {
			return errors.New(fmt.Sprintf("instances[%d] is required", index))
		}

		instance.ID = strings.TrimSpace(instance.ID)
		instance.DisplayName = strings.TrimSpace(instance.DisplayName)
		instance.BaseURL = strings.TrimSpace(instance.BaseURL)
		instance.TokenFile = strings.TrimSpace(instance.TokenFile)
		if instance.ID == "" {
			return errors.New(fmt.Sprintf("instances[%d].id is required", index))
		}
		if _, ok := seen[instance.ID]; ok {
			return errors.New(fmt.Sprintf("duplicate instance id %q", instance.ID))
		}
		seen[instance.ID] = struct{}{}

		if instance.DisplayName == "" {
			instance.DisplayName = instance.ID
		}
		if instance.BaseURL == "" {
			return errors.New(fmt.Sprintf("instances[%d].base_url is required", index))
		}
		if instance.TokenFile == "" {
			return errors.New(fmt.Sprintf("instances[%d].token_file is required", index))
		}
		if instance.TokenFile != "" && !filepath.IsAbs(instance.TokenFile) {
			instance.TokenFile = filepath.Clean(filepath.Join(configDir, instance.TokenFile))
		}
	}

	for index, instance := range cfg.Instances {
		for previousIndex := 0; previousIndex < index; previousIndex++ {
			previous := cfg.Instances[previousIndex]
			if strings.EqualFold(instance.DisplayName, previous.DisplayName) {
				return errors.New(fmt.Sprintf(
					"duplicate instance display_name %q for instances %q and %q",
					instance.DisplayName,
					previous.ID,
					instance.ID,
				))
			}
		}
		for otherIndex, other := range cfg.Instances {
			if otherIndex != index && strings.EqualFold(instance.DisplayName, other.ID) {
				return errors.New(fmt.Sprintf(
					"instance display_name %q for %q conflicts with instance id %q",
					instance.DisplayName,
					instance.ID,
					other.ID,
				))
			}
		}
	}

	return nil
}

// GetServer returns the listener configuration.
func GetServer() *ServerConfig {
	if config == nil {
		return nil
	}
	return config.Server
}

// GetLogger returns the logging configuration.
func GetLogger() *LoggerConfig {
	if config == nil {
		return nil
	}
	return config.Logger
}

// GetInstances returns all configured SafeLine instances.
func GetInstances() []*InstanceConfig {
	if config == nil {
		return nil
	}
	return config.Instances
}
