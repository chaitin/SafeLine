package config

type Config struct {
	Log          LogConfig
	MgtWebserver string
}

func DefaultGlobalConfig() Config {
	return Config{
		Log:          DefaultLogConfig(),
		MgtWebserver: "169.254.0.4:9002",
	}
}
