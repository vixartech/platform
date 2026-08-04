ALLOY_SERVICE ?= alloy

ALLOY_IMAGE  ?= v1.11.0
ALLOY_CONFIG ?= $(PLATFORM_DIR)/deployments/alloy/configs/config.alloy

export ALLOY_IMAGE
export ALLOY_CONFIG

ALLOY_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/alloy/compose.yaml
PLATFORM_COMPOSE   += -f $(ALLOY_COMPOSE_FILE)

.PHONY: infra/alloy/config
infra/alloy/config:
	$(PLATFORM_COMPOSE) config $(ALLOY_SERVICE)

.PHONY: infra/alloy/up
infra/alloy/up:
	$(PLATFORM_COMPOSE) up -d $(ALLOY_SERVICE)

.PHONY: infra/alloy/start
infra/alloy/start:
	$(PLATFORM_COMPOSE) start $(ALLOY_SERVICE)

.PHONY: infra/alloy/stop
infra/alloy/stop:
	$(PLATFORM_COMPOSE) stop $(ALLOY_SERVICE)

.PHONY: infra/alloy/restart
infra/alloy/restart:
	$(PLATFORM_COMPOSE) restart $(ALLOY_SERVICE)

.PHONY: infra/alloy/ps
infra/alloy/ps:
	$(PLATFORM_COMPOSE) ps $(ALLOY_SERVICE)

.PHONY: infra/alloy/logs
infra/alloy/logs:
	$(PLATFORM_COMPOSE) logs -f $(ALLOY_SERVICE)

.PHONY: infra/alloy/shell
infra/alloy/shell:
	$(PLATFORM_COMPOSE) exec $(ALLOY_SERVICE) sh

.PHONY: infra/alloy/remove
infra/alloy/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(ALLOY_SERVICE)
