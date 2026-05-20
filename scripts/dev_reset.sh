#!/usr/bin/env bash
set -e
docker compose -f deploy/docker-compose.yml down -v --remove-orphans
docker compose -f deploy/docker-compose.yml up -d postgres
sleep 3
bash scripts/migrate_all.sh
bash scripts/seed_data.sh
