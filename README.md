# xk6-ocsp
A [k6](https://k6.io) extension to test [OCSP](https://datatracker.ietf.org/doc/html/rfc6960) responders.

Feel free to send PRs, current limitations:
- no support for "exotic" ECC curves (e.g. brainpool) in certificates
- RSASSAPSS is not supported for OCSP signatures

To workaround these limitations see the documentation in the examples directory.

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:
  ```bash
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```bash
  $ xk6 build --with github.com/nikosn/xk6-ocsp@latest --output k6-ocsp-check
  ```
in case of problems try
  ```bash
  $ GOWORK=off xk6 build --with github.com/nikosn/xk6-ocsp@latest --output k6-ocsp-check
  ```

## Usage
Check the examples directory for sample k6 scripts.

To import the ocsp module
```JavaScript
import ocspmodule from 'k6/x/ocsp';
```

### ExtractSerialNumberAndOCSPURIFromCert
```go
ocspmodule.ExtractSerialNumberAndOCSPURIFromCert(certPath string) (string, string, error)
```
ExtractSerialNumberAndOCSPURIFromCert extracts the serialNumber and OCSP URI from a PEM encoded certificate
The serialNumber is returned as HEX string. This does not work with "exotic" ECC keys like brainpool.

### CreateRequest
```go
ocspmodule.CreateRequest(hexSerialNumber string, issuerCertPath string, hashAlgorithm string) ([]byte, string, error)
```
CreateOCSPRequest creates an OCSP request using the given hex serialNumber and issuer certificate path where the PEM encoded issuer certificate is placed into.
this does not work with "exotic" ECC keys like brainpool  
hashAlgorithm can be SHA1 or SHA256.

### CheckResponse
```go
ocspmodule.CheckResponse(ocspResponseBytes []byte, verifySignature bool) (string, error)
```
CheckOCSPResponse checks the OCSP response. Signature verification fails in case custom ECC curves like brainpool are used. RSAPSS signatures aren't supported either.  
To workaround this set verifySignature to false.
