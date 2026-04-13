#!/bin/sh

set -e

echo "run db migration"

migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "starting the app"

exec "$@"