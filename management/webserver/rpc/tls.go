package rpc

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
)

const (
	// TLSEnv names the environment variable that decides how the control
	// channel is secured.
	//
	//   - "0" / off / false / no: serve in the clear (warned).
	//   - anything else, including unset: TLS is required. Credentials come
	//     from `mgt -gen_certs`. Missing files are an error; there is no
	//     silent fallback to plaintext (the subscription token would otherwise
	//     travel on the same network as the containers).
	TLSEnv = "MGT_GRPC_TLS"

	// CertDirEnv names the environment variable that holds the directory of
	// the server credentials.
	CertDirEnv = "MGT_GRPC_CERT_DIR"

	// ClientCertDirEnv names the environment variable that holds the directory
	// the credentials of tcontrollerd are published in.
	ClientCertDirEnv = "MGT_GRPC_CLIENT_CERT_DIR"

	// ServerCertFile, ServerKeyFile and CertAuthorityFile are the names of the
	// server credentials inside CertDir().
	ServerCertFile    = "grpc_server.crt"
	ServerKeyFile     = "grpc_server.key"
	CertAuthorityFile = "grpc_ca.crt"

	// ClientCertFile and ClientKeyFile are the names of the credentials of
	// tcontrollerd inside ClientCertDir().
	ClientCertFile = "tcd_client.crt"
	ClientKeyFile  = "tcd_client.key"

	// ServerName is the name the server certificate is issued for.
	//
	// tcontrollerd reaches the management server at an address of the
	// deployment (169.254.0.4:9002 in the standard compose file), which no
	// certificate can be issued for. The name is therefore fixed on both sides
	// and checked against the certificate, which is what makes the peer
	// verifiable: the address may change, the name may not.
	ServerName = "waf-management-server"
)

// CertDir returns the directory the server credentials of the control channel
// are stored in.
func CertDir() string {
	if dir := strings.TrimSpace(os.Getenv(CertDirEnv)); dir != "" {
		return dir
	}

	return filepath.Join(config.GlobalConfig.MgtResDir, constants.CertsPath)
}

// ClientCertDir returns the directory the credentials of tcontrollerd are
// published in.
//
// It is a directory below the one both containers mount, because the tengine
// container is where tcontrollerd runs and it does not mount the resources
// directory of the management container.
func ClientCertDir() string {
	if dir := strings.TrimSpace(os.Getenv(ClientCertDirEnv)); dir != "" {
		return dir
	}

	return filepath.Join(filepath.Dir(TokenFileFromEnv()), "grpc")
}

// ServerCreds returns the transport credentials of the control channel.
//
// A nil option means the channel is served in the clear, and that only happens
// when TLSEnv is an explicit off value. Unset TLSEnv still requires TLS: the
// deployment is expected to have run `mgt -gen_certs`; missing credentials are
// returned as an error instead of a plaintext listener.
func ServerCreds() (grpc.ServerOption, error) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(TLSEnv))) {
	case "0", "off", "false", "no":
		logger.Warnf("%s is disabled: the control channel is served in the clear, and the subscription token is the only gate in front of it. Anything that shares the network of this container can read it.", TLSEnv)
		return nil, nil
	}

	creds, err := loadServerCreds()
	if err != nil {
		return nil, err
	}

	logger.Infof("The control channel is served over TLS and only accepts a client certificate issued by %s", filepath.Join(CertDir(), CertAuthorityFile))

	return grpc.Creds(creds), nil
}

// loadServerCreds reads the server credentials and builds the TLS configuration
// of the control channel from them.
func loadServerCreds() (credentials.TransportCredentials, error) {
	certPath := filepath.Join(CertDir(), ServerCertFile)
	keyPath := filepath.Join(CertDir(), ServerKeyFile)

	certificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "load %s", certPath)
	}

	clientCAs, err := certPool(filepath.Join(CertDir(), CertAuthorityFile))
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		// The subscription token already keeps a caller that only reaches the
		// port out, but it travels over this connection: without a client
		// certificate, anyone who can read the traffic of the deployment can
		// replay it and subscribe as tcontrollerd.
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  clientCAs,
		MinVersion: tls.VersionTLS12,
	}), nil
}

// certPool reads a PEM encoded certificate authority.
func certPool(path string) (*x509.CertPool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read %s", path)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(content) {
		return nil, errors.New("no certificate found in " + path)
	}

	return pool, nil
}
