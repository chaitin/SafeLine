package config

import (
	"chaitin.cn/dev/go/settings"
)

type GRPCConfig struct {
	ListenAddr string `yaml:"listen_addr"`
}

func DefaultGRPCConfig() GRPCConfig {
	return GRPCConfig{
		ListenAddr: ":9002",
	}
}

func (sc *GRPCConfig) Load(setting *settings.Setting) error {
	if err := setting.Unmarshal("grpc_server", sc); err != nil {
		return err
	}

	return nil
}
