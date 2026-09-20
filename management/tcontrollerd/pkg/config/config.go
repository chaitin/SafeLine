package config

import (
	"os"
	"path/filepath"
	"strings"

	"chaitin.cn/dev/go/settings"
)

var (
	GlobalConfig = DefaultGlobalConfig()
)

// DefaultMgtTokenFile is where the control channel token is read from by
// default. The standard deployment mounts the same host directory at
// /app/sock inside the management container, which creates the file, and
// inside the tengine container, which is where tcontrollerd runs.
const DefaultMgtTokenFile = "/app/sock/mgt_grpc_token"

// CertDirEnv names the environment variable that holds the directory of the
// credentials the management server published for this side of the control
// channel. The default is below the directory of the token file, which is the
// directory the standard deployment shares between the two containers.
const CertDirEnv = "TCD_MGT_CERT_DIR"

// CertDir returns the directory the credentials of the control channel are
// read from.
func CertDir() string {
	if v, ok := os.LookupEnv(CertDirEnv); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}

	return filepath.Join(filepath.Dir(GlobalConfig.MgtTokenFile), "grpc")
}

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

	// The token file is optional in the configuration file: the default is
	// where the management server puts it.
	if err := s.Unmarshal("mgt_token_file", &GlobalConfig.MgtTokenFile); err != nil {
		return err
	}

	if v, ok := os.LookupEnv("MGT_TOKEN_FILE"); ok {
		GlobalConfig.MgtTokenFile = v
	}

	return nil
}
