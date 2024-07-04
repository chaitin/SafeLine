package cmd

import (
	"crypto/x509/pkix"
	"path/filepath"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

func GenCerts() error {
	if err := genServerCert(); err != nil {
		return err
	}
	if err := genClientCACert(); err != nil {
		return err
	}
	return nil
}

func genServerCert() error {
	certPath := filepath.Join(config.GlobalConfig.MgtResDir, constants.CertsPath, "server.crt")
	keyPath := filepath.Join(config.GlobalConfig.MgtResDir, constants.CertsPath, "server.key")
	if err := utils.WriteCertIfNotExist(
		certPath,
		keyPath,
		func() ([]byte, []byte, error) {
			return utils.GenerateCert(
				[]string{},
				3650,
				4096,
				&pkix.Name{
					Country:            []string{"CN"},
					Province:           []string{"Beijing"},
					Locality:           []string{"Beijing"},
					Organization:       []string{"Beijing WAF Technology Co., Ltd."},
					OrganizationalUnit: []string{"Service Infrastructure Department"},
					CommonName:         "WAF Management Server",
				},
				false,
			)
		}); err != nil {
		return err
	}
	return nil
}

func genClientCACert() error {
	certPath := filepath.Join(config.GlobalConfig.MgtResDir, constants.CertsPath, "client_ca.crt")
	keyPath := filepath.Join(config.GlobalConfig.MgtResDir, constants.CertsPath, "client_ca.key")
	if err := utils.WriteCertIfNotExist(
		certPath,
		keyPath,
		func() ([]byte, []byte, error) {
			return utils.GenerateCert(
				[]string{},
				3650,
				4096,
				&pkix.Name{
					Country:            []string{"CN"},
					Province:           []string{"Beijing"},
					Locality:           []string{"Beijing"},
					Organization:       []string{"Beijing WAF Technology Co., Ltd."},
					OrganizationalUnit: []string{"Service Infrastructure Department"},
					CommonName:         "WAF Client Certificate Authority",
				},
				true,
			)
		}); err != nil {
		return err
	}
	return nil
}
