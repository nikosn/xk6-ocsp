# check using an openssl created OCSP request without signature verification.
This can be used to test "exotic" OCSP responders, e.g. using brainpool curves or RSASSAPSS.
The provided ocsp.js needs the actual OCSP payload (URL) as environment variables.
The OCSP responder will be called via http get.

Just use
```bash
$ export ENDPOINT_URL=$(./get_req.sh /path/to/cert_to_check.pem /path/to/issuer_ca.pem)
$ k6-ocsp-check run ocsp.js
```
# check by creating an OCSP request in each iteration including signature verification.
The script ocsp-full.js needs following environment variables:
```bash
$ export CERT_PATH=/path/to/cert_to_check.pem
$ export ISSUER_CERT_PATH=/path/to/issuer_ca.pem
$ export HASH_ALGORITHM=SHA1
```
Instead of SHA1, SHA256 can also be defined as HASH_ALGORITHM.  
The implementation will try to parse the OCSP AIA extension from the to be checked certificate.
```bash
$ k6-ocsp-check run ocsp-full.js
```

# check by creating an OCSP request from a list of certificate serialnumbers in each iteration including signature verification.
The script ocsp-with-serialNumbers-full.js needs following environment variables:
```bash
$ export ENDPOINT_URL=http://...
$ export HEX_SERIALNUMBERS_FILE_PATH=/path/to/line_separated_hex_serials.txt
$ export ISSUER_CERT_PATH=/path/to/issuer_ca.pem
$ export HASH_ALGORITHM=SHA1
```
Instead of SHA1, SHA256 can also be defined as HASH_ALGORITHM.  
The implementation will use a random selected serialnumber from the provided in each iteration to create the OCSP request.
```bash
$ k6-ocsp-check run ocsp-full.js
```