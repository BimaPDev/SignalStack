#!/bin/bash
set -e

# Load .env from project root
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../.env"

if [ -f "$ENV_FILE" ]; then
    export $(grep -v '^#' "$ENV_FILE" | xargs)
fi

if [ -z "$POSTGRES_ADDR" ]; then
    echo "ERROR: POSTGRES_ADDR is not set"
    exit 1
fi

echo "Seeding database..."
psql "$POSTGRES_ADDR" -f "$SCRIPT_DIR/../migrations/seed.sql"
echo "Done."
