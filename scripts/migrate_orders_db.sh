#!/bin/sh
set -e

echo "Run ${DB_ORDER_DATABASE} migration..."
migrate -path=/migrations \
  -database="postgres://${DB_USERNAME}:${DB_PASSWORD}@postgres:5432/${DB_ORDER_DATABASE}?sslmode=disable" \
  up
