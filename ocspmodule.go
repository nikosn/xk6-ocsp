package ocspmodule

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"

	"github.com/nikosn/xk6-ocsp/ocsp"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/ocsp", new(Ocspmodule))
}

type Ocspmodule struct {
}

// ExtractSerialNumberAndOCSPURIFromCert extracts the serialNumber and OCSP URI from a PEM encoded certificate
// the serialNumber is returned as HEX string
// this does not work with "exotic" ECC keys like brainpool
func (o *Ocspmodule) ExtractSerialNumberAndOCSPURIFromCert(certPath string) (string, string, error) {
	var ocspAIA = ""
	var serialNumber = ""
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return serialNumber, ocspAIA, fmt.Errorf("failed to read certificate: %w", err)
	}
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return serialNumber, ocspAIA, fmt.Errorf("failed to decode PEM block containing certificate")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return serialNumber, ocspAIA, fmt.Errorf("failed to parse certificate: %w", err)
	}
	if len(cert.OCSPServer) > 0 {
		ocspAIA = cert.OCSPServer[0]
	}
	if ocspAIA == "" {
		return serialNumber, ocspAIA, fmt.Errorf("failed to get OCSP uri from certificate", err)
	}
	return cert.SerialNumber.Text(16), ocspAIA, nil
}

// CreateOCSPRequest creates an OCSP request using the given hex serialNumber and issuer certificate path where the PEM encoded issuer certificate is placed into.
// this does not work with "exotic" ECC keys like brainpool
func (o *Ocspmodule) CreateRequest(hexSerialNumber, issuerCertPath, hashAlgorithm string) ([]byte, error) {

	issuerCertPEM, err := os.ReadFile(issuerCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read issuer certificate: %w", err)
	}

	issuerCertBlock, _ := pem.Decode(issuerCertPEM)
	if issuerCertBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing issuer certificate")
	}

	issuerCert, err := x509.ParseCertificate(issuerCertBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issuer certificate: %w", err)
	}

	// hexSerialNumber to big.init
	serialNumber := new(big.Int)
	serialNumber.SetString(hexSerialNumber, 16)

	// Determine hash algorithm
	var hash crypto.Hash
	switch hashAlgorithm {
	case "SHA256":
		hash = crypto.SHA256
	case "SHA1":
		hash = crypto.SHA1
	default:
		return nil, fmt.Errorf("unsupported hash algorithm: %s", hashAlgorithm)
	}

	// Create OCSP request
	ocspRequest, err := ocsp.CreateRequest(serialNumber, issuerCert, &ocsp.RequestOptions{
		Hash: hash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OCSP request: %w", err)
	}

	return ocspRequest, nil
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
