#!/bin/bash

# Define the directory where the certificate will reside
PROJECT_DIR="/home/tng/dev/api-oa-integrator"
CERT_DIR="$PROJECT_DIR/cert"
PRIVATE_KEY="$CERT_DIR/private-key.pem"
CERTIFICATE="$CERT_DIR/certificate.pem"
CSR="$CERT_DIR/csr.pem"

rm -rf $CERT_DIR
mkdir -p $CERT_DIR
# Generate a new private key and certificate
openssl genpkey -algorithm RSA -out $PRIVATE_KEY
openssl req -new -x509 -days 365 -key $PRIVATE_KEY -out $CERTIFICATE -subj "/C=MY/ST=KL/L=KL/O=Enctech Services Sdn Bhd/OU=OA/CN=OA/emailAddress=zamri@enctechgroup.com"
openssl req -new -key $PRIVATE_KEY -out $CSR -subj "/C=MY/ST=KL/L=KL/O=Enctech Services Sdn Bhd/OU=OA/CN=OA/emailAddress=zamri@enctechgroup.com"

cd $PROJECT_DIR || true

# Copy cert to all directories
make copy_cert

# Turn down the containers
docker compose down
docker compose up -d --build