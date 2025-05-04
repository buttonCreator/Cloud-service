#!/bin/sh
set -e

until pg_isready -h $POSTGRES_HOST -p 5432 -U $POSTGRES_USER; do
  echo "Waiting for PostgreSQL to start..."
  sleep 1
done

/usr/local/bin/migrate -path /migrations -database "$POSTGRESQL_CONNECTION_STRING" -verbose up

exec /cloud-service
