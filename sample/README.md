# check using an openssl created OCSP request without signature verification

The provided ocsp.js needs the actual OCSP payload (URL) as environment variables.
The OCSP responder will be called via http get.

Just use
```
export ENDPOINT_URL=$(./get_req.sh /path/to/cert_to_check.pem /path/to/issuer_ca.pem)
k6-ocsp-check run ocsp.js
```
# check by creating an OCSP request in each iteration including signature verification.
The script needs following environment variables:
```
export CERT_PATH=/path/to/cert_to_check.pem
export ISSUER_CERT_PATH=/path/to/issuer_ca.pem
export HASH_ALGORITHM=SHA1
```
The OCSP responder URL will try to parse the OCSP AIA extension from the certificate to check.
```
k6-ocsp-check run ocsp-full.js
```