package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"
)

var (
	notFoundError   = errors.New("Certificate is not found.")
	notStartedError = errors.New("Certificate valid date is not started.")
	shortLifeError  = errors.New("Certificate will be invalid after given duration.")
)

// TLSHealthChecker implements health.Checkable
type TLSHealthChecker struct {
	cert *x509.Certificate
	dur  time.Duration
}

// Healthy returns error if it could not loaded certificate or validity bounds invalid
func (hc *TLSHealthChecker) Healthy() error {
	if hc.cert == nil {
		return notFoundError
	}
	if hc.cert.NotBefore.After(time.Now()) {
		return notStartedError
	}
	if hc.cert.NotAfter.Before(time.Now().Add(hc.dur)) {
		return shortLifeError
	}
	return nil
}

// New returns TLSHealthChecker
func New(certFile, keyFile string, d time.Duration) *TLSHealthChecker {
	c, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return &TLSHealthChecker{cert: nil}
	}
	if len(c.Certificate) > 0 {
		c1, _ := x509.ParseCertificate(c.Certificate[0])
		return &TLSHealthChecker{cert: c1, dur: d}
	}
	return &TLSHealthChecker{cert: nil}
}

// NewWithCert returns TLSHealthChecker
func NewWithCert(cert *x509.Certificate, d time.Duration) *TLSHealthChecker {
	return &TLSHealthChecker{cert: cert, dur: d}
}
