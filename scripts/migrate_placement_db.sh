#!/bin/sh
set -e

echo "Run  migration placement_db..."
migrate -path=/migrations \
  -database="postgres://test:test@postgres:5432/placement_db?sslmode=disable" \
  up
