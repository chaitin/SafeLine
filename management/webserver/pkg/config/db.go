package config

import (
	"net/url"
	"os"

	"chaitin.cn/dev/go/settings"
)

type DBConfig struct {
	URL     string `yaml:"url"`
	LogSQL  bool   `yaml:"log_sql"`
	SSLMode bool   `yaml:"ssl_mode"`
}

func DefaultDBConfig() DBConfig {
	return DBConfig{
		URL:     "postgres://safeline-ce:safeline-ce@127.0.0.1/safeline-ce",
		LogSQL:  false,
		SSLMode: false,
	}
}

func (dbc *DBConfig) Load(setting *settings.Setting) error {
	if err := setting.Unmarshal("db", dbc); err != nil {
		return err
	}

	if v, ok := os.LookupEnv("DATABASE_URL"); ok {
		dbc.URL = v
	}

	dbURL, err := url.Parse(dbc.URL)
	if err != nil {
		return err
	}
	q := dbURL.Query()
	if !dbc.SSLMode {
		q.Set("sslmode", "disable")
	}
	dbURL.RawQuery = q.Encode()
	dbc.URL = dbURL.String()

	return nil
}
