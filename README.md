# xk6-ocsp
k6 extension to test OCSP responders

Feel free to send PRs, as this does not support brainpool curves at all and RSASSAPSS for OCSP signatures.

# Setup
Just clone this repository and
1. `go install go.k6.io/xk6/cmd/xk6@latest`
2. `xk6 build --with xk6-ocsp=. --output k6-ocsp-check`

# Usage
Check the sample directory for a k6 script.