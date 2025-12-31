#!/bin/sh
set -e

echo "Run migration ${DB_INVENTORY_DATABASE}..."
migrate -path=/migrations \
  -database="postgres://${DB_USERNAME}:${DB_PASSWORD}@postgres:5432/${DB_INVENTORY_DATABASE}?sslmode=disable" \
  up
