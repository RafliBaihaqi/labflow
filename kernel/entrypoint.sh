#!/bin/sh
set -e

# Validate required env vars fail-fast
: "${DATABASE_URL:?DATABASE_URL required}"
: "${MIGRATIONS_URL:?MIGRATIONS_URL required}"
: "${REDIS_URL:?REDIS_URL required}"
: "${JWT_SECRET:?JWT_SECRET required}"
: "${COOKIE_SECRET:?COOKIE_SECRET required}"
: "${MINIO_ENDPOINT:?MINIO_ENDPOINT required}"
: "${MINIO_ACCESS:?MINIO_ACCESS required}"
: "${MINIO_SECRET:?MINIO_SECRET required}"

echo "=== LabFlow Go Backend Entrypoint ==="
echo "Checking migrations dir..."
ls -la /app/sql/migrations | head -20

# Retry migrations up to 10 times (postgres might still be starting)
MAX_RETRIES=10
RETRY=0
until [ $RETRY -ge $MAX_RETRIES ]; do
  echo "Running database migrations (attempt $((RETRY+1))/$MAX_RETRIES)..."
  if migrate -path /app/sql/migrations -database "${MIGRATIONS_URL}" up; then
    echo "Migrations complete."
    break
  fi
  RETRY=$((RETRY+1))
  echo "Migration failed, retrying in 3s... ($RETRY/$MAX_RETRIES)"
  sleep 3
done

if [ $RETRY -ge $MAX_RETRIES ]; then
  echo "Migrations failed after $MAX_RETRIES attempts" >&2
  exit 1
fi

echo "Starting server..."
exec ./main "$@"
