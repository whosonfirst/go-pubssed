package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"path"
	"time"
)

func EnsureTLSRoot() (string, error) {

	tmp := os.TempDir()
	pony := path.Join(tmp, "httpony")
	root := path.Join(pony, "certificates")

	info, err := os.Stat(root)

	if os.IsNotExist(err) {
		err = os.MkdirAll(root, 0700)

		if err != nil {
			return "", nil
		}
	} else {

		if info.Mode().Perm() != 0700 {
			return "", errors.New("Certificate root has insecure permissions")
		}
	}

	return root, nil
}

/*
	the following is basically just this: https://github.com/mattrobenolt/https
	still to do: custom root; custom TTL and org details; filenames
*/

func GenerateTLSCert(host string, root string) (string, string, error) {

	info, err := os.Stat(root)

	if os.IsNotExist(err) {
		return "", "", errors.New("Certificate root does not exist!")
	}

	if info.Mode().Perm() != 0700 {
		return "", "", errors.New("Certificate root has insecure permissions")
	}

	root = path.Join(root, host)

	cert_path := path.Join(root, "cert.pem")
	key_path := path.Join(root, "key.pem")

	cert_exists := false
	key_exists := false

	_, err = os.Stat(cert_path)

	if !os.IsNotExist(err) {
		cert_exists = true
	}

	_, err = os.Stat(key_path)

	if !os.IsNotExist(err) {
		key_exists = true
	}

	if cert_exists && key_exists {
		return cert_path, key_path, nil
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return "", "", err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		IsCA: true,

		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)

	if err != nil {
		return "", "", err
	}

	err = os.MkdirAll(root, 0700)

	if err != nil {
		return "", "", err
	}

	certOut, err := os.Create(cert_path)

	if err != nil {
		return "", "", err
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(key_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return "", "", err
	}

	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	return cert_path, key_path, nil
}
