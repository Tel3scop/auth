#!/bin/bash
source .env

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "host=${MIGRATION_DIR} port=${MIGRATION_DIR} dbname=${MIGRATION_DIR} user=${MIGRATION_DIR} password=${MIGRATION_DIR} sslmode=${MIGRATION_DIR}" up -v