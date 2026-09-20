package cmd

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"
	"path/filepath"
	"time"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/rpc"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

func GenCerts() error {
	if err := genServerCert(); err != nil {
		return err
	}
	if err := genClientCACert(); err != nil {
		return err
	}
	if err := genControlChannelCerts(); err != nil {
		return err
	}
	return nil
}

// genControlChannelCerts issues the credentials of the control channel between
// the management server and tcontrollerd.
//
// The channel carries the complete site configuration of the installation, and
// it used to be served in the clear with a token in the metadata as its only
// gate: the token is readable by anything that shares the network of the
// containers, and a caller that reads it can subscribe as tcontrollerd. Both
// ends are therefore issued by a CA that only this installation knows, and each
// end checks the other against it.
//
// The server credentials stay in the resources directory, which only the
// management container mounts. The credentials of tcontrollerd are published in
// the directory both containers mount, next to the token file, because
// tcontrollerd runs in the tengine container.
func genControlChannelCerts() error {
	caCert, caKey, err := controlChannelCA()
	if err != nil {
		return err
	}

	// The name is fixed on both sides: tcontrollerd dials an address of the
	// deployment, which no certificate can be issued for. The addresses of the
	// standard deployment are listed as well, so that a client which verifies
	// the address instead of the name still works.
	hostnames := []string{rpc.ServerName, "mgt", "safeline-mgt", "localhost"}

	if err := ensureSignedCert(
		filepath.Join(rpc.CertDir(), rpc.ServerCertFile),
		filepath.Join(rpc.CertDir(), rpc.ServerKeyFile),
		caCert,
		caKey,
		hostnames,
		&pkix.Name{
			Country:            []string{"CN"},
			Province:           []string{"Beijing"},
			Locality:           []string{"Beijing"},
			Organization:       []string{"Beijing WAF Technology Co., Ltd."},
			OrganizationalUnit: []string{"Service Infrastructure Department"},
			CommonName:         rpc.ServerName,
		},
		false,
	); err != nil {
		return err
	}

	clientDir := rpc.ClientCertDir()

	if err := ensureSignedCert(
		filepath.Join(clientDir, rpc.ClientCertFile),
		filepath.Join(clientDir, rpc.ClientKeyFile),
		caCert,
		caKey,
		[]string{},
		&pkix.Name{
			Country:            []string{"CN"},
			Province:           []string{"Beijing"},
			Locality:           []string{"Beijing"},
			Organization:       []string{"Beijing WAF Technology Co., Ltd."},
			OrganizationalUnit: []string{"Service Infrastructure Department"},
			CommonName:         "tcontrollerd",
		},
		true,
	); err != nil {
		return err
	}

	// tcontrollerd verifies the server with the CA, so it needs a copy. A
	// certificate is public, unlike the key that never leaves the management
	// container.
	return utils.EnsureRenameWriteFile(filepath.Join(clientDir, rpc.CertAuthorityFile), caCert, 0644)
}

// controlChannelCA reads the certificate authority of the control channel,
// creating it when it is missing.
func controlChannelCA() ([]byte, []byte, error) {
	certPath := filepath.Join(rpc.CertDir(), rpc.CertAuthorityFile)
	keyPath := filepath.Join(rpc.CertDir(), "grpc_ca.key")

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
					CommonName:         "WAF Control Channel Certificate Authority",
				},
				true,
			)
		}); err != nil {
		return nil, nil, err
	}

	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}

// ensureSignedCert writes a certificate signed by the CA of the control
// channel, computing it again when the one on disk does not belong to that CA.
//
// Both ends check the other against the CA, so a certificate that was issued by
// an older CA, or that is about to expire, only shows up as a handshake that
// fails for no visible reason. It is replaced here instead; the key of the
// management server credentials is kept to the owner, the one of tcontrollerd
// is read by the tengine container and is kept to the owner as well.
func ensureSignedCert(certPath, keyPath string, caCert, caKey []byte, hostnames []string, subject *pkix.Name, client bool) error {
	reusable, err := signedByAuthority(certPath, caCert)
	if err != nil {
		return err
	}

	if reusable {
		if _, err := os.Stat(keyPath); err == nil {
			return nil
		}
	}

	cert, key, err := utils.GenerateSignedCert(caCert, caKey, hostnames, 3650, 4096, subject, client)
	if err != nil {
		return err
	}

	if err := utils.EnsureRenameWriteFile(certPath, cert, 0644); err != nil {
		return err
	}

	return utils.EnsureRenameWriteFile(keyPath, key, 0600)
}

// signedByAuthority reports whether the certificate at path is signed by the
// given authority and is not about to expire.
func signedByAuthority(path string, caCertPEM []byte) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return false, nil
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, nil
	}

	caBlock, _ := pem.Decode(caCertPEM)
	if caBlock == nil {
		return false, nil
	}

	authority, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return false, nil
	}

	if err := certificate.CheckSignatureFrom(authority); err != nil {
		return false, nil
	}

	return time.Now().Add(24 * time.Hour).Before(certificate.NotAfter), nil
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
