LOCAL_BIN:=$(CURDIR)/bin

.PHONY: help
help: ## List all available targets with help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: install-golangci-lint
install-golangci-lint: ## install linter
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

.PHONY: lint
lint: ## run linter
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml

.PHONY: install-deps
install-deps: ## install dependencies
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
.PHONY: get-deps
get-deps: ## get dependencies
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
.PHONY: generate
generate: ## generate handlers
	make generate-user-api

.PHONY: generate-user-api
generate-user-api:  ## generate handlers for /api/user
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user_v1/user.proto

local-migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	goose -dir ${MIGRATION_DIR} postgres  "host=${POSTGRES_HOST} port=${POSTGRES_PORT} dbname=${POSTGRES_DB} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} sslmode=${POSTGRES_SSLMODE}" up -v

local-migration-down:
	goose -dir ${MIGRATION_DIR} postgres  "host=${POSTGRES_HOST} port=${POSTGRES_PORT} dbname=${POSTGRES_DB} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} sslmode=${POSTGRES_SSLMODE}" down -v

create-migration:
	if [ -z "$(name)" ]; then \
			echo "Укажите название миграции в формате 'make create-migration name=create_users'"; \
		else \
			goose -dir ./migrations create $(name) sql; \
		fi

up:
	docker-compose up -d
down:
	docker-compose down