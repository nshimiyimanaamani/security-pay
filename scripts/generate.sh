#!/bin/bash

echo "generating rsa keys..."
openssl genrsa -out secrets/jwt.private.pem 2048
openssl rsa -in secrets/jwt.private.pem -pubout > secrets/jwt.public.pem

echo "exporting keys..."
export PRIVKEY=$(cat secrets/jwt.private.pem)
export PUBKEY=$(cat secrets/jwt.public.pem)