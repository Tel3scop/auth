#!/bin/bash
source migration.env

sleep 5 && goose -dir "${MIGRATION_DIR}" postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} dbname=${POSTGRES_DB} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} sslmode=${POSTGRES_SSLMODE}" up -v