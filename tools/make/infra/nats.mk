NATS_SERVICE ?= nats

NATS_IMAGE        ?= 2.12.2-alpine
NATS_PORT         ?= 4222
NATS_MONITOR_PORT ?= 8222

export NATS_IMAGE
export NATS_PORT
export NATS_MONITOR_PORT

NATS_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/nats/compose.yaml
PLATFORM_COMPOSE  += -f $(NATS_COMPOSE_FILE)

.PHONY: infra/nats/config
infra/nats/config:
	$(PLATFORM_COMPOSE) config $(NATS_SERVICE)

.PHONY: infra/nats/up
infra/nats/up:
	$(PLATFORM_COMPOSE) up -d $(NATS_SERVICE)

.PHONY: infra/nats/start
infra/nats/start:
	$(PLATFORM_COMPOSE) start $(NATS_SERVICE)

.PHONY: infra/nats/stop
infra/nats/stop:
	$(PLATFORM_COMPOSE) stop $(NATS_SERVICE)

.PHONY: infra/nats/restart
infra/nats/restart:
	$(PLATFORM_COMPOSE) restart $(NATS_SERVICE)

.PHONY: infra/nats/ps
infra/nats/ps:
	$(PLATFORM_COMPOSE) ps $(NATS_SERVICE)

.PHONY: infra/nats/logs
infra/nats/logs:
	$(PLATFORM_COMPOSE) logs -f $(NATS_SERVICE)

.PHONY: infra/nats/shell
infra/nats/shell:
	$(PLATFORM_COMPOSE) exec $(NATS_SERVICE) sh

.PHONY: infra/nats/remove
infra/nats/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(NATS_SERVICE)
