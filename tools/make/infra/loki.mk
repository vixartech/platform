LOKI_SERVICE ?= loki

LOKI_IMAGE  ?= 3.6.1
LOKI_CONFIG ?= $(PLATFORM_DIR)/deployments/loki/configs/loki.yaml

export LOKI_IMAGE
export LOKI_CONFIG

LOKI_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/loki/compose.yaml
PLATFORM_COMPOSE  += -f $(LOKI_COMPOSE_FILE)

.PHONY: infra/loki/config
infra/loki/config:
	$(PLATFORM_COMPOSE) config $(LOKI_SERVICE)

.PHONY: infra/loki/up
infra/loki/up:
	$(PLATFORM_COMPOSE) up -d $(LOKI_SERVICE)

.PHONY: infra/loki/start
infra/loki/start:
	$(PLATFORM_COMPOSE) start $(LOKI_SERVICE)

.PHONY: infra/loki/stop
infra/loki/stop:
	$(PLATFORM_COMPOSE) stop $(LOKI_SERVICE)

.PHONY: infra/loki/restart
infra/loki/restart:
	$(PLATFORM_COMPOSE) restart $(LOKI_SERVICE)

.PHONY: infra/loki/ps
infra/loki/ps:
	$(PLATFORM_COMPOSE) ps $(LOKI_SERVICE)

.PHONY: infra/loki/logs
infra/loki/logs:
	$(PLATFORM_COMPOSE) logs -f $(LOKI_SERVICE)

.PHONY: infra/loki/shell
infra/loki/shell:
	$(PLATFORM_COMPOSE) exec $(LOKI_SERVICE) sh

.PHONY: infra/loki/remove
infra/loki/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(LOKI_SERVICE)
