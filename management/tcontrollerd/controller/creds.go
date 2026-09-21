package controller

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/credentials"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/config"
)

const (
	// TLSEnv names the environment variable that used to allow a plaintext
	// control channel. Plaintext is no longer supported: a value of "0" /
	// off / false / no is an error.
	TLSEnv = "TCD_MGT_TLS"

	// ServerCertFile is the certificate the management server is expected to
	// present, ClientCertFile and ClientKeyFile are the credentials issued for
	// this side of the control channel, and CertAuthorityFile is the authority
	// both of them are issued by.
	ServerName        = "waf-management-server"
	ClientCertFile    = "tcd_client.crt"
	ClientKeyFile     = "tcd_client.key"
	CertAuthorityFile = "grpc_ca.crt"
)

// transportCreds returns the credentials of the control channel.
//
// The channel is always reached over TLS. Missing credentials or an explicit
// request for plaintext are errors.
func transportCreds() (credentials.TransportCredentials, error) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(TLSEnv))) {
	case "0", "off", "false", "no":
		return nil, errors.New("plaintext is not supported: unset " + TLSEnv + " and use credentials published by mgt -gen_certs")
	}

	return loadTransportCreds()
}

// loadTransportCreds reads the credentials the management server published for
// this side of the control channel.
func loadTransportCreds() (credentials.TransportCredentials, error) {
	dir := config.CertDir()

	certPath := filepath.Join(dir, ClientCertFile)
	certificate, err := tls.LoadX509KeyPair(certPath, filepath.Join(dir, ClientKeyFile))
	if err != nil {
		return nil, errors.Wrapf(err, "load %s", certPath)
	}

	caPath := filepath.Join(dir, CertAuthorityFile)
	content, err := os.ReadFile(caPath)
	if err != nil {
		return nil, errors.Wrapf(err, "read %s", caPath)
	}

	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(content) {
		return nil, errors.New("no certificate found in " + caPath)
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      roots,
		// The server is reached at an address of the deployment, which no
		// certificate is issued for: the name it is issued for is fixed on
		// both sides instead, and it is checked against the certificate.
		ServerName: ServerName,
		MinVersion: tls.VersionTLS12,
	}), nil
}
