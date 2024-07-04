package config

import (
	"chaitin.cn/dev/go/settings"
)

type TelemetryConfig struct {
	Addr string `yaml:"addr"`
}

func DefaultTelemetryConfig() TelemetryConfig {
	return TelemetryConfig{
		Addr: "rivers-telemetry:10086",
	}
}

func (t *TelemetryConfig) Load(setting *settings.Setting) error {
	if err := setting.Unmarshal("telemetry", t); err != nil {
		return err
	}

	return nil
}
