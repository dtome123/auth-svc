#!/bin/sh

set -e

mkdir -p cert

if [ ! -f cert/private.pem ]; then
  echo "🔐 Generating cert..."
  ./bin/genrsa --out-dir=cert
else
  echo "✅ Cert already exists. Skipping."
fi


echo "🔑 Public key:"
cat cert/public.pem

echo "🚀 Starting main app..."
exec ./main
