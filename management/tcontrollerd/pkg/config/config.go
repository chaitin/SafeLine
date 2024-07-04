package config

import (
	"os"

	"chaitin.cn/dev/go/settings"
)

var (
	GlobalConfig = DefaultGlobalConfig()
)

func InitConfigs(configFilePath string) error {
	s, err := settings.New(configFilePath)
	if err != nil {
		return err
	}

	if err = GlobalConfig.Log.Load(s); err != nil {
		return err
	}

	if err := s.Unmarshal("mgt_addr", &GlobalConfig.MgtWebserver); err != nil {
		return err
	}

	if v, ok := os.LookupEnv("MGT_ADDR"); ok {
		GlobalConfig.MgtWebserver = v
	}

	return nil
}
