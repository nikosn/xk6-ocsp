# xk6-ocsp
k6 extension to test OCSP responders

go install go.k6.io/xk6/cmd/xk6@latest

xk6 build --with xk6-ocsp=. --output k6-ocsp-check
