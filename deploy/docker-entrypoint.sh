#!/bin/sh
set -e

mkdir -p /app/uploads
chown -R appuser:appuser /app/uploads /app/migrations 2>/dev/null || true

exec su-exec appuser "$@"
