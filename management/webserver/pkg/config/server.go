package config

import (
	"chaitin.cn/dev/go/settings"
)

type ServerConfig struct {
	ListenAddr  string `yaml:"listen_addr"`
	DevMode     bool   `yaml:"dev_mode"`
	IntenseMode bool   `yaml:"intense_mode"`
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		ListenAddr: ":9001",
		DevMode:    false,
	}
}

func (sc *ServerConfig) Load(setting *settings.Setting) error {
	if err := setting.Unmarshal("server", sc); err != nil {
		return err
	}

	return nil
}
