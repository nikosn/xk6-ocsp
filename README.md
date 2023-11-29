# xk6-ocsp
k6 extension to test OCSP responders

Feel free to send PRs, as this does not support brainpool curves at all and RSASSAPSS for OCSP signatures.

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
Check the sample directory for sample k6 scripts.