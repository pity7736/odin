#!/usr/bin/env bash
#
# odin-cli demo: exercises the full client flow against a running server.
# Start the server first:  go run ./src
# Then run:                ./scripts/odin-cli-demo.sh
#
set -euo pipefail

cd "$(dirname "$0")/.."

EMAIL="demo@odin.dev"
PASSWORD="demo1234"
CLI="go run ./cmd/odin-cli"

echo "==> register"
$CLI register -email "$EMAIL" -password "$PASSWORD"

echo "==> login"
$CLI login -email "$EMAIL" -password "$PASSWORD"

echo "==> create chunks (example finance data, encrypted client-side)"
$CLI create-chunk -email "$EMAIL" \
  -plaintext '{"type":"account","name":"Checking","currency":"USD","balance":1000}'
$CLI create-chunk -email "$EMAIL" \
  -plaintext '{"type":"account","name":"Savings","currency":"USD","balance":5000}'
$CLI create-chunk -email "$EMAIL" \
  -plaintext '{"type":"category","name":"Food"}'
$CLI create-chunk -email "$EMAIL" \
  -plaintext '{"type":"income","amount":3000,"currency":"USD","description":"Salary"}'
$CLI create-chunk -email "$EMAIL" \
  -plaintext '{"type":"expense","amount":25.5,"currency":"USD","category":"Food","description":"Lunch"}'

echo "==> done"
