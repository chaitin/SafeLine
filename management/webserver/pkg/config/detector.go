package config

import (
	"chaitin.cn/dev/go/settings"
)

type DetectorConfig struct {
	Addr        string `yaml:"addr"`
	FslBytecode string `yaml:"fsl_bytecode"`
}

func DefaultDetectorConfig() DetectorConfig {
	return DetectorConfig{
		Addr:        "http://127.0.0.1:8001",
		FslBytecode: "bytecode",
	}
}

func (d *DetectorConfig) Load(setting *settings.Setting) error {
	if err := setting.Unmarshal("detector", d); err != nil {
		return err
	}

	return nil
}
