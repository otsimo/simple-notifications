package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math"
	"math/big"
	"net"
	"os"
	"testing"
	"time"
)

//Cert creation code created from https://github.com/coreos/pkg/blob/master/k8s-tlsutil/k8s-tlsutil.go
const (
	RSAKeySize   = 1024
	Duration365d = time.Hour * 24 * 365
)

type CertConfig struct {
	CommonName   string
	Organization []string
	AltNames     AltNames
}

// AltNames contains the domain names and IP addresses that will be added
// to the API Server's x509 certificate SubAltNames field. The values will
// be passed directly to the x509.Certificate object.
type AltNames struct {
	DNSNames []string
	IPs      []net.IP
}

func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, RSAKeySize)
}
func EncodePrivateKeyPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}

func EncodeCertificatePEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}
func NewSelfSignedCACertificate(cfg CertConfig, key *rsa.PrivateKey, startDate time.Time, validDuration time.Duration) (*x509.Certificate, error) {
	dur := Duration365d * 10
	if validDuration != 0 {
		dur = validDuration
	}

	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		NotBefore:             startDate,
		NotAfter:              startDate.Add(dur),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA: true,
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

func ParsePEMEncodedCACert(pemdata []byte) (*x509.Certificate, error) {
	decoded, _ := pem.Decode(pemdata)
	if decoded == nil {
		return nil, errors.New("no PEM data found")
	}
	return x509.ParseCertificate(decoded.Bytes)
}

func ParsePEMEncodedPrivateKey(pemdata []byte) (*rsa.PrivateKey, error) {
	decoded, _ := pem.Decode(pemdata)
	if decoded == nil {
		return nil, errors.New("no PEM data found")
	}
	return x509.ParsePKCS1PrivateKey(decoded.Bytes)
}

func NewSignedCertificate(cfg CertConfig, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey, startTime time.Time, validDuration time.Duration) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}

	dur := Duration365d
	if validDuration != 0 {
		dur = validDuration
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: caCert.Subject.Organization,
		},
		DNSNames:     cfg.AltNames.DNSNames,
		IPAddresses:  cfg.AltNames.IPs,
		SerialNumber: serial,
		NotBefore:    caCert.NotBefore,
		NotAfter:     startTime.Add(dur),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

func NewTestCert(startTime time.Time, validDuration time.Duration) (*x509.Certificate, *rsa.PrivateKey, error) {
	cfg := CertConfig{
		CommonName:   "test",
		Organization: []string{"Otsimo"},
		AltNames:     AltNames{},
	}
	pkca, err := NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	ca, err := NewSelfSignedCACertificate(cfg, pkca, startTime, Duration365d)
	if err != nil {
		return nil, nil, err
	}
	pk, err := NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	sc, err := NewSignedCertificate(cfg, pk, ca, pkca, startTime, validDuration)
	if err != nil {
		return nil, nil, err
	}
	return sc, pk, nil
}

func CreateTestCert(dir string, startTime time.Time, validDuration time.Duration) (string, string, error) {
	c, pk, err := NewTestCert(startTime, validDuration)
	if err != nil {
		return "", "", err
	}
	crtfile, err := ioutil.TempFile(dir, "test-cert.pem")
	if err != nil {
		return "", "", err
	}
	if _, err := crtfile.Write(EncodeCertificatePEM(c)); err != nil {
		return "", "", err
	}
	if err := crtfile.Close(); err != nil {
		return "", "", err
	}
	keyfile, err := ioutil.TempFile(dir, "test-key.pem")
	if err != nil {
		return "", "", err
	}
	if _, err := keyfile.Write(EncodePrivateKeyPEM(pk)); err != nil {
		return "", "", err
	}
	if err := keyfile.Close(); err != nil {
		return "", "", err
	}
	return crtfile.Name(), keyfile.Name(), nil
}

func TestInvalidCert(t *testing.T) {
	checker := NewWithCert(nil, time.Hour*24*14)
	if err := checker.Healthy(); err != notFoundError {
		t.Fatalf("want='%v' got='%v'", notFoundError, err)
	}
}

func TestNotBefore(t *testing.T) {
	d := time.Hour * 24
	now := time.Now()
	start := now.Add(d)
	c, _, err := NewTestCert(start, d)
	if err != nil {
		t.Fatal(err)
	}
	checker := NewWithCert(c, time.Hour*6)
	if err := checker.Healthy(); err != notStartedError {
		t.Fatalf("want='%v' got='%v'", notStartedError, err)
	}
}

func TestNotAfter(t *testing.T) {
	d := time.Hour * 24 * 7 //one week
	now := time.Now()
	c, _, err := NewTestCert(now, d)
	if err != nil {
		t.Fatal(err)
	}
	checker := NewWithCert(c, time.Hour*24*8)
	if err := checker.Healthy(); err != shortLifeError {
		t.Fatalf("want='%v' got='%v'", shortLifeError, err)
	}
	checker2 := NewWithCert(c, time.Hour*24*3)
	if err := checker2.Healthy(); err != nil {
		t.Fatalf("want='%v' got='%v'", nil, err)
	}
}

func TestReadingCertificate(t *testing.T) {
	d := time.Hour * 24 * 7 //one week
	now := time.Now()
	c, k, err := CreateTestCert("", now, d)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(c) // clean up cert
	defer os.Remove(k) // clean up key

	checker := New(c, k, time.Hour*24*8)
	if err := checker.Healthy(); err != shortLifeError {
		t.Fatalf("want='%v' got='%v'", shortLifeError, err)
	}
	checker2 := New(c, k, time.Hour*24*3)
	if err := checker2.Healthy(); err != nil {
		t.Fatalf("want='%v' got='%v'", nil, err)
	}
	checker3 := New("asd", "asdasf", time.Hour*24*3)
	if err := checker3.Healthy(); err != notFoundError {
		t.Fatalf("want='%v' got='%v'", notFoundError, err)
	}
}
