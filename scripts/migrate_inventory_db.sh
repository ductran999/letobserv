#!/bin/sh
set -e

echo "Run migration inventory_db..."
migrate -path=/migrations \
  -database="postgres://test:test@postgres:5432/inventory_db?sslmode=disable" \
  up
