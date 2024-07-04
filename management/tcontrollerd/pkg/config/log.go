package config

import (
	"chaitin.cn/dev/go/settings"
)

type LogConfig struct {
	Output string `yaml:"output"`
	Level  string `yaml:"level"`
}

func DefaultLogConfig() LogConfig {
	return LogConfig{
		Output: "stdout",
		Level:  "info",
	}
}

func (lc *LogConfig) Load(setting *settings.Setting) error {
	if err := setting.Unmarshal("log", lc); err != nil {
		return err
	}

	return nil
}
