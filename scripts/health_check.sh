#!/usr/bin/env bash
set -e
for p in 8080 8081 8082 8083 8084 8085 8086 8087 8088; do curl -fsS http://localhost:$p/health >/dev/null; done
echo health ok
