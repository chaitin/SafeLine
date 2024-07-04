package config

type Config struct {
	Log          LogConfig
	DB           DBConfig
	Server       ServerConfig
	Detector     DetectorConfig
	Telemetry    TelemetryConfig
	GPRC         GRPCConfig
	PlatformAddr string
	MgtResDir    string
	NgxResDir    string
}

func DefaultGlobalConfig() Config {
	return Config{
		Log:          DefaultLogConfig(),
		DB:           DefaultDBConfig(),
		Server:       DefaultServerConfig(),
		Detector:     DefaultDetectorConfig(),
		Telemetry:    DefaultTelemetryConfig(),
		GPRC:         DefaultGRPCConfig(),
		PlatformAddr: "https://waf-ce.chaitin.cn",
		MgtResDir:    "/resources/management",
		NgxResDir:    "/resources/nginx",
	}
}
