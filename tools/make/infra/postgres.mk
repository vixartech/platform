POSTGRES_SERVICE ?= postgres

POSTGRES_IMAGE    ?= 18.4-alpine3.24
POSTGRES_DB       ?= db
POSTGRES_USER     ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_PORT     ?= 5432
POSTGRES_CONFIG   ?= $(PLATFORM_DIR)/deployments/postgres/configs/postgresql.conf

export POSTGRES_IMAGE
export POSTGRES_DB
export POSTGRES_USER
export POSTGRES_PASSWORD
export POSTGRES_PORT
export POSTGRES_CONFIG

POSTGRES_MIGRATIONS ?= $(CURDIR)/migrations
POSTGRES_DSN        ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
POSTGRES_MIGRATIONS_TABLE ?= $(notdir $(CURDIR))_schema_migrations
POSTGRES_MIGRATION_DSN    ?= $(POSTGRES_DSN)&x-migrations-table=$(POSTGRES_MIGRATIONS_TABLE)

POSTGRES_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/postgres/compose.yaml
PLATFORM_COMPOSE      += -f $(POSTGRES_COMPOSE_FILE)

.PHONY: infra/postgres/config
infra/postgres/config:
	$(PLATFORM_COMPOSE) config $(POSTGRES_SERVICE)

.PHONY: infra/postgres/up
infra/postgres/up:
	$(PLATFORM_COMPOSE) up -d $(POSTGRES_SERVICE)

.PHONY: infra/postgres/start
infra/postgres/start:
	$(PLATFORM_COMPOSE) start $(POSTGRES_SERVICE)

.PHONY: infra/postgres/stop
infra/postgres/stop:
	$(PLATFORM_COMPOSE) stop $(POSTGRES_SERVICE)

.PHONY: infra/postgres/restart
infra/postgres/restart:
	$(PLATFORM_COMPOSE) restart $(POSTGRES_SERVICE)

.PHONY: infra/postgres/ps
infra/postgres/ps:
	$(PLATFORM_COMPOSE) ps $(POSTGRES_SERVICE)

.PHONY: infra/postgres/logs
infra/postgres/logs:
	$(PLATFORM_COMPOSE) logs -f $(POSTGRES_SERVICE)

.PHONY: infra/postgres/shell
infra/postgres/shell:
	$(PLATFORM_COMPOSE) exec $(POSTGRES_SERVICE) psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)"

.PHONY: infra/postgres/remove
infra/postgres/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(POSTGRES_SERVICE)

.PHONY: infra/postgres/migration
infra/postgres/migrations: bin/migrate
	@$(BIN_PATH)/migrate create -ext sql -dir "$(POSTGRES_MIGRATIONS)" migration

.PHONY: infra/postgres/migrate
infra/postgres/migrate: bin/migrate
	@$(BIN_PATH)/migrate -path "$(POSTGRES_MIGRATIONS)" -database "$(POSTGRES_MIGRATION_DSN)" up

.PHONY: infra/postgres/rollback
infra/postgres/rollback: bin/migrate
	@$(BIN_PATH)/migrate -path "$(POSTGRES_MIGRATIONS)" -database "$(POSTGRES_MIGRATION_DSN)" down 1
