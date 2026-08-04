GRAFANA_SERVICE ?= grafana

GRAFANA_IMAGE        ?= 12.1.0
GRAFANA_PORT         ?= 3000
GRAFANA_PROVISIONING ?= $(PLATFORM_DIR)/deployments/grafana/provisioning

export GRAFANA_IMAGE
export GRAFANA_PORT
export GRAFANA_PROVISIONING

GRAFANA_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/grafana/compose.yaml
PLATFORM_COMPOSE     += -f $(GRAFANA_COMPOSE_FILE)

.PHONY: infra/grafana/config
infra/grafana/config:
	$(PLATFORM_COMPOSE) config $(GRAFANA_SERVICE)

.PHONY: infra/grafana/up
infra/grafana/up:
	$(PLATFORM_COMPOSE) up -d $(GRAFANA_SERVICE)

.PHONY: infra/grafana/start
infra/grafana/start:
	$(PLATFORM_COMPOSE) start $(GRAFANA_SERVICE)

.PHONY: infra/grafana/stop
infra/grafana/stop:
	$(PLATFORM_COMPOSE) stop $(GRAFANA_SERVICE)

.PHONY: infra/grafana/restart
infra/grafana/restart:
	$(PLATFORM_COMPOSE) restart $(GRAFANA_SERVICE)

.PHONY: infra/grafana/ps
infra/grafana/ps:
	$(PLATFORM_COMPOSE) ps $(GRAFANA_SERVICE)

.PHONY: infra/grafana/logs
infra/grafana/logs:
	$(PLATFORM_COMPOSE) logs -f $(GRAFANA_SERVICE)

.PHONY: infra/grafana/shell
infra/grafana/shell:
	$(PLATFORM_COMPOSE) exec $(GRAFANA_SERVICE) sh

.PHONY: infra/grafana/remove
infra/grafana/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(GRAFANA_SERVICE)
