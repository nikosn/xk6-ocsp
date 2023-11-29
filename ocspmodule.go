package ocspmodule

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/nikosn/xk6-ocsp/ocsp"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/ocsp", new(Ocspmodule))
}

type Ocspmodule struct {
}

// CreateOCSPRequest creates an OCSP request using the given certificate and issuer certificate paths where the PEM encoded certs are placed into.
// this does not work with "exotic" ECC keys like brainpool
func (o *Ocspmodule) CreateRequest(certPath, issuerCertPath, hashAlgorithm string) ([]byte, string, error) {
	var ocspAIA = ""
	// Load certificate and issuer certificate
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, ocspAIA, fmt.Errorf("failed to read certificate: %w", err)
	}

	issuerCertPEM, err := os.ReadFile(issuerCertPath)
	if err != nil {
		return nil, ocspAIA, fmt.Errorf("failed to read issuer certificate: %w", err)
	}

	// Parse certificate and issuer certificate
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, ocspAIA, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, ocspAIA, fmt.Errorf("failed to parse certificate: %w", err)
	}

	issuerCertBlock, _ := pem.Decode(issuerCertPEM)
	if issuerCertBlock == nil {
		return nil, ocspAIA, fmt.Errorf("failed to decode PEM block containing issuer certificate")
	}

	issuerCert, err := x509.ParseCertificate(issuerCertBlock.Bytes)
	if err != nil {
		return nil, ocspAIA, fmt.Errorf("failed to parse issuer certificate: %w", err)
	}

	// Determine hash algorithm
	var hash crypto.Hash
	switch hashAlgorithm {
	case "SHA256":
		hash = crypto.SHA256
	case "SHA1":
		hash = crypto.SHA1
	default:
		return nil, ocspAIA, fmt.Errorf("unsupported hash algorithm: %s", hashAlgorithm)
	}

	// Create OCSP request
	ocspRequest, err := ocsp.CreateRequest(cert, issuerCert, &ocsp.RequestOptions{
		Hash: hash,
	})
	if err != nil {
		return nil, ocspAIA, fmt.Errorf("failed to create OCSP request: %w", err)
	}

	// extrac OCSP uri from cert
	if len(cert.OCSPServer) > 0 {
		ocspAIA = cert.OCSPServer[0]
	}
	if ocspAIA == "" {
		return nil, ocspAIA, fmt.Errorf("failed to get OCSP uri from certificate", err)
	}

	return ocspRequest, ocspAIA, nil
}

// CheckOCSPResponse checks the OCSP response.
// signature verification fails in case custom ECC curves like brainpool are used. RSAPSS signatures aren't supported either.
// to workaround this set verifySignature to false
func (o *Ocspmodule) CheckResponse(ocspResponseBytes []byte, verifySignature bool) (string, error) {
	ocspResponse, err := ocsp.ParseResponse(ocspResponseBytes, nil, verifySignature)
	if err != nil {
		return "", fmt.Errorf("failed to parse OCSP response: %w", err)
	}

	// Check the validity of the OCSP response
	status := ocspResponse.Status
	switch status {
	case ocsp.Good:
		return "Good", nil
	case ocsp.Revoked:
		return "Revoked", nil
	case ocsp.Unknown:
		return "Unknown", nil
	default:
		return fmt.Sprintf("OCSP error status code: %v", status), nil
	}
}
