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
	// TLSEnv names the environment variable that decides how the control
	// channel is reached. "1" requires TLS and refuses to start without it,
	// "0" reaches the server in the clear, and anything else (including an
	// unset value) uses TLS as soon as the management server published the
	// credentials, which it does when its container runs `mgt -gen_certs` on
	// start.
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
// A nil result means the channel is reached in the clear. The management
// server serves it over TLS and only accepts a client certificate it issued,
// so a deployment in which the two containers were not upgraded together keeps
// working through the switch below, and every other deployment fails the
// handshake rather than subscribing without a certificate.
func transportCreds() (credentials.TransportCredentials, error) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(TLSEnv))) {
	case "0", "off", "false", "no":
		logger.Warnf("%s is disabled: the control channel is reached in the clear, and the subscription token is readable by anything that shares the network of this container.", TLSEnv)
		return nil, nil
	}

	creds, err := loadTransportCreds()
	if err != nil {
		if required := strings.ToLower(strings.TrimSpace(os.Getenv(TLSEnv))); required == "1" || required == "on" || required == "true" || required == "yes" {
			return nil, err
		}

		logger.Warnf("Reaching the control channel in the clear: %s. Set %s=1 to make this an error.", err, TLSEnv)
		return nil, nil
	}

	return creds, nil
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
