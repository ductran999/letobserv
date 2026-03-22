#!/bin/sh

# inventory seeds
shopt -s nullglob

for f in ./services/inventory/migrations/seeds/*.sql; do
  echo "Seeding $f"
  psql "postgres://test:test@localhost:5432/inventory_db?sslmode=disable" -f "$f"
done
