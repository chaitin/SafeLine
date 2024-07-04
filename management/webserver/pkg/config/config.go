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

	if err = GlobalConfig.DB.Load(s); err != nil {
		return err
	}

	if err = GlobalConfig.Log.Load(s); err != nil {
		return err
	}

	if err = GlobalConfig.Server.Load(s); err != nil {
		return err
	}

	if err = GlobalConfig.Detector.Load(s); err != nil {
		return err
	}

	if err = GlobalConfig.Telemetry.Load(s); err != nil {
		return err
	}

	if err = GlobalConfig.GPRC.Load(s); err != nil {
		return err
	}

	if err := settings.Unmarshal("platform_addr", &GlobalConfig.PlatformAddr); err != nil {
		return err
	}

	if v, ok := os.LookupEnv("MANAGEMENT_RESOURCES_DIR"); ok {
		GlobalConfig.MgtResDir = v
	}

	if v, ok := os.LookupEnv("NGINX_RESOURCES_DIR"); ok {
		GlobalConfig.NgxResDir = v
	}

	return nil
}
