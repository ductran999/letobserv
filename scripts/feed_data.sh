#!/bin/sh

# inventory seeds
shopt -s nullglob

for f in ./migrations/inventory/seeds/*.sql; do
  echo "Seeding $f"
  psql "postgres://${DB_USERNAME}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_INVENTORY_DATABASE}?sslmode=disable" -f "$f"
done
