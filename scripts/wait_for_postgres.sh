#!/usr/bin/env bash
until pg_isready -d ${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/hypercompute?sslmode=disable}; do sleep 1; done
