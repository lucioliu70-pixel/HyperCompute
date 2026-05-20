#!/usr/bin/env bash
set -e
command -v docker >/dev/null; docker compose version >/dev/null; command -v go >/dev/null; command -v python3 >/dev/null
if [ ! -f .env ]; then cp .env.example .env; fi
docker compose -f deploy/docker-compose.yml up -d postgres redis nats prometheus grafana
sleep 5
scripts/migrate_all.sh
docker compose -f deploy/docker-compose.yml up -d
sleep 10
scripts/seed_data.sh
scripts/health_check.sh
scripts/integration_test.sh
echo Gateway URL: http://localhost:8080
echo Admin URL: http://localhost:3000
echo Grafana URL: http://localhost:3001
echo Default API Key: hc_live_test_xxx
