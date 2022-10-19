package chaser

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"errors"
)

func (c *Chaser) ProvisionCertificate() {
	keyExists := true
	certExists := true

	if _, err := os.Stat("chaser.key"); errors.Is(err, os.ErrNotExist) {
		keyExists = false
	}
	if _, err := os.Stat("chaser.crt"); errors.Is(err, os.ErrNotExist) {
		certExists = false
	}

	if keyExists && !certExists {
		fmt.Printf("ERROR: Invalid Configuration\n")
		fmt.Printf("Found chaser.key but missing chaser.crt\n")
		fmt.Printf("Delete chaser.key OR restore chaser.crt\n")
		panic("invalid configuration")
	}
	if !keyExists && certExists {
		fmt.Printf("ERROR: Invalid Configuration\n")
		fmt.Printf("Found chaser.crt but missing chaser.key\n")
		fmt.Printf("Delete chaser.crt OR restore chaser.key\n")
		panic("invalid configuration")
	}
	if keyExists && certExists {
		c.ReadCertificate()
	}

	if !keyExists && !certExists {
		c.ProvisionNewCertificate()
		c.ReadCertificate()
	}

	fmt.Printf("To add this chaser your local snitch, run 'precrux remote add %s'\n\n", c.Name)
	fmt.Printf("Then copy+paste the following text into your precrux snitch\n\n")
	fmt.Printf("%s\n\n\n", c.certBytes)
}
func (c *Chaser) ReadCertificate() {
	var err error
	c.certBytes, err = os.ReadFile("chaser.crt")
	if err != nil {
		panic(err)
	}
	c.keyBytes, err = os.ReadFile("chaser.key")
	if err != nil {
		panic(err)
	}
}

func (c *Chaser) ProvisionNewCertificate() {
	fmt.Printf("Generating chaser certificate, please wait...\n\n")

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	// PEM encoding of private key
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
	)
	// fmt.Println(string(keyPEM))
	if err = os.WriteFile("chaser.key", keyPEM, 0600); err != nil {
		panic(err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * 100 * time.Hour)

	//Certificate template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(0),
		Subject:               pkix.Name{CommonName: c.Name}, // legacy
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:              []string{c.Name},
	}

	//Create certificate using template
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	//pem encoding of certificate
	certPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		},
	)
	// fmt.Println(string(certPem))
	if err = os.WriteFile("chaser.crt", certPem, 0600); err != nil {
		panic(err)
	}
}
