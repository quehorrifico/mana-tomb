#!/bin/bash
set -e

echo "🔄 Checking if database exists..."
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'mana_tomb'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE mana_tomb"

echo "✅ Database ensured: mana_tomb"