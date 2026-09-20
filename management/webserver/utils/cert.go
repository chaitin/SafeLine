package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"chaitin.cn/dev/go/errors"
)

func GenerateCert(hostnames []string, days int64, keyBits int, subject *pkix.Name, isCA bool) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	duration := int64(time.Hour) * 24 * days
	notAfter := notBefore.Add(time.Duration(duration))

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	var keyUsage x509.KeyUsage
	if isCA {
		keyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	} else {
		keyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		DNSNames:              hostnames,
		BasicConstraintsValid: true,
		IsCA:                  isCA,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage:              keyUsage,
	}
	if subject != nil {
		template.Subject = *subject
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	var certBuffer bytes.Buffer
	err = pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, nil, err
	}

	var keyBuffer bytes.Buffer
	err = pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return nil, nil, err
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}

func WriteCertIfNotExist(certFilePath, keyFilePath string, generator func() ([]byte, []byte, error)) error {
	cert, key, err := genCertIfNotExist(certFilePath, keyFilePath, generator)
	if err != nil {
		return err
	}

	if err = EnsureRenameWriteFile(certFilePath, cert, 0644); err != nil {
		return err
	}

	// A certificate is public, its private key is not: 0644 lets every local
	// user read the key and impersonate the service it belongs to.
	if err = EnsureRenameWriteFile(keyFilePath, key, 0600); err != nil {
		return err
	}

	// The file may already have existed with the permissive mode that older
	// releases gave it, in which case the mode passed to the write above is
	// not applied by the kernel.
	if err = os.Chmod(keyFilePath, 0600); err != nil {
		return err
	}

	return nil
}

func genCertIfNotExist(certFilePath, keyFilePath string, generator func() ([]byte, []byte, error)) ([]byte, []byte, error) {
	exist, err := FilesExist(certFilePath, keyFilePath)
	if err != nil {
		return nil, nil, err
	}

	if !exist {
		return generator()
	} else {
		cert, err := ioutil.ReadFile(certFilePath)
		if err != nil {
			return nil, nil, err
		}

		key, err := ioutil.ReadFile(keyFilePath)
		if err != nil {
			return nil, nil, err
		}
		return cert, key, nil
	}
}

// GenerateSignedCert issues a certificate for hostnames that is signed by the
// certificate authority that certPEM and keyPEM belong to.
//
// GenerateCert produces a self signed certificate, which is what the console
// uses: the browser is told to accept it. The control channel cannot work that
// way, because both ends have to recognise each other without trusting the
// network they talk over, so the management server issues the credentials of
// both ends from a CA that only this installation knows.
func GenerateSignedCert(certPEM, keyPEM []byte, hostnames []string, days int64, keyBits int, subject *pkix.Name, client bool) ([]byte, []byte, error) {
	caCert, caKey, err := parseCertAuthority(certPEM, keyPEM)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	if client {
		extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	notBefore := time.Now()
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(time.Duration(days) * 24 * time.Hour),
		DNSNames:              hostnames,
		BasicConstraintsValid: true,
		IsCA:                  false,
		ExtKeyUsage:           extKeyUsage,
		KeyUsage:              keyUsage,
	}
	if subject != nil {
		template.Subject = *subject
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &privateKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	var certBuffer bytes.Buffer
	if err = pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, err
	}

	var keyBuffer bytes.Buffer
	if err = pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return nil, nil, err
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}

// parseCertAuthority reads the certificate and the private key of a CA.
func parseCertAuthority(certPEM, keyPEM []byte) (*x509.Certificate, *rsa.PrivateKey, error) {
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, nil, errors.New("the certificate authority is not PEM encoded")
	}

	caCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, nil, errors.New("the private key of the certificate authority is not PEM encoded")
	}

	caKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return caCert, caKey, nil
}
