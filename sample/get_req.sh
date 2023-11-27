#!/bin/bash

# use ./get_env.sh /path/to/cert_to_check.pem /path/to/cert_ca_cert.pem

cert=$1
issuer=$2

req=$(openssl ocsp -issuer $issuer -cert $cert -reqout - | base64 -w 0)
uri=$(openssl x509 -in $cert -noout -ocsp_uri)

echo "$uri/$req"
