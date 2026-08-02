#!/bin/bash

# Define the directory where the certificate will reside
PROJECT_DIR="/home/tng/dev/api-oa-integrator"
CERT_DIR="$PROJECT_DIR/cert"
PRIVATE_KEY="$CERT_DIR/private-key.pem"
CERTIFICATE="$CERT_DIR/certificate.pem"

rm -rf $CERT_DIR
mkdir -p $CERT_DIR
# Generate a new private key and certificate
openssl genpkey -algorithm RSA -out $PRIVATE_KEY -pkeyopt rsa_keygen_bits:2048
openssl req -new -x509 -days 365 -key $PRIVATE_KEY -out $CERTIFICATE -subj "/C=US/ST=State/L=City/O=Organization/OU=Unit/CN=example.com"

cd $PROJECT_DIR || true

# Copy cert to all directories
make copy_cert

# Turn down the containers
docker compose down
docker compose up -d --build