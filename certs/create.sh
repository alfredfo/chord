#!/bin/bash

cat > cert.v3.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = myserver.local
DNS.2 = myserver1.local
IP.1 = 192.168.1.1
IP.2 = 192.168.2.1
IP.3 = 127.0.0.1
EOF

openssl x509 -req -in cert.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out signed.crt -days 730 -sha256 -extfile cert.v3.ext

cat signed.crt ca.crt > combined.crt
