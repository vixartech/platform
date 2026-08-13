LISTMONK_SERVICE ?= listmonk

LISTMONK_IMAGE          ?= v6.2.0
LISTMONK_PORT           ?= 9000
LISTMONK_ADMIN_USER     ?= admin
LISTMONK_ADMIN_PASSWORD ?= password
LISTMONK_DB             ?= listmonk

export LISTMONK_IMAGE
export LISTMONK_PORT
export LISTMONK_ADMIN_USER
export LISTMONK_ADMIN_PASSWORD
export LISTMONK_DB

LISTMONK_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/listmonk/compose.yaml
PLATFORM_COMPOSE      += -f $(LISTMONK_COMPOSE_FILE)

.PHONY: infra/listmonk/config
infra/listmonk/config:
	$(PLATFORM_COMPOSE) config $(LISTMONK_SERVICE)

.PHONY: infra/listmonk/up
infra/listmonk/up: infra/listmonk/database
	$(PLATFORM_COMPOSE) up -d $(LISTMONK_SERVICE)

.PHONY: infra/listmonk/database
infra/listmonk/database:
	-$(PLATFORM_COMPOSE) exec -T $(POSTGRES_SERVICE) createdb -U "$$POSTGRES_USER" "$(LISTMONK_DB)"

.PHONY: infra/listmonk/install
infra/listmonk/install: infra/listmonk/database
	$(PLATFORM_COMPOSE) run --rm --no-deps $(LISTMONK_SERVICE) ./listmonk --install --idempotent --yes

.PHONY: infra/listmonk/seed
infra/listmonk/seed: infra/listmonk/database
	$(PLATFORM_COMPOSE) exec -T $(POSTGRES_SERVICE) psql -U "$$POSTGRES_USER" -d "$(LISTMONK_DB)" -v ON_ERROR_STOP=1 < "$(PLATFORM_DIR)/deployments/listmonk/sql/seed.sql"

.PHONY: infra/listmonk/setup
infra/listmonk/setup: infra/mailpit/up infra/listmonk/install infra/listmonk/up infra/listmonk/seed infra/listmonk/restart

.PHONY: infra/listmonk/upgrade
infra/listmonk/upgrade:
	$(PLATFORM_COMPOSE) run --rm --no-deps $(LISTMONK_SERVICE) ./listmonk --upgrade --yes

.PHONY: infra/listmonk/start
infra/listmonk/start:
	$(PLATFORM_COMPOSE) start $(LISTMONK_SERVICE)

.PHONY: infra/listmonk/stop
infra/listmonk/stop:
	$(PLATFORM_COMPOSE) stop $(LISTMONK_SERVICE)

.PHONY: infra/listmonk/restart
infra/listmonk/restart:
	$(PLATFORM_COMPOSE) restart $(LISTMONK_SERVICE)

.PHONY: infra/listmonk/ps
infra/listmonk/ps:
	$(PLATFORM_COMPOSE) ps $(LISTMONK_SERVICE)

.PHONY: infra/listmonk/logs
infra/listmonk/logs:
	$(PLATFORM_COMPOSE) logs -f $(LISTMONK_SERVICE)

.PHONY: infra/listmonk/shell
infra/listmonk/shell:
	$(PLATFORM_COMPOSE) exec $(LISTMONK_SERVICE) sh

.PHONY: infra/listmonk/remove
infra/listmonk/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(LISTMONK_SERVICE)
