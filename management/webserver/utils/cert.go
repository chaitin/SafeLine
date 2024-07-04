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
	"time"
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

	if err = EnsureRenameWriteFile(keyFilePath, key, 0644); err != nil {
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
