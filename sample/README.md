The provided ocsp.js needs the actual OCSP payload (URL) as environment variables.

Just use
```
export ENDPOINT_URL=$(./get_req.sh /path/to/cert_to_check.pem /path/to/issuer_ca.pem)
xk6-ocsp/k6-ocsp-check run ocsp.js
```
