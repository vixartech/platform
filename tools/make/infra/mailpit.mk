MAILPIT_SERVICE ?= mailpit

MAILPIT_IMAGE     ?= v1.30.7
MAILPIT_SMTP_PORT ?= 1025
MAILPIT_HTTP_PORT ?= 8025

export MAILPIT_IMAGE
export MAILPIT_SMTP_PORT
export MAILPIT_HTTP_PORT

MAILPIT_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/mailpit/compose.yaml
PLATFORM_COMPOSE     += -f $(MAILPIT_COMPOSE_FILE)

.PHONY: infra/mailpit/config
infra/mailpit/config:
	$(PLATFORM_COMPOSE) config $(MAILPIT_SERVICE)

.PHONY: infra/mailpit/up
infra/mailpit/up:
	$(PLATFORM_COMPOSE) up -d $(MAILPIT_SERVICE)

.PHONY: infra/mailpit/start
infra/mailpit/start:
	$(PLATFORM_COMPOSE) start $(MAILPIT_SERVICE)

.PHONY: infra/mailpit/stop
infra/mailpit/stop:
	$(PLATFORM_COMPOSE) stop $(MAILPIT_SERVICE)

.PHONY: infra/mailpit/restart
infra/mailpit/restart:
	$(PLATFORM_COMPOSE) restart $(MAILPIT_SERVICE)

.PHONY: infra/mailpit/ps
infra/mailpit/ps:
	$(PLATFORM_COMPOSE) ps $(MAILPIT_SERVICE)

.PHONY: infra/mailpit/logs
infra/mailpit/logs:
	$(PLATFORM_COMPOSE) logs -f $(MAILPIT_SERVICE)

.PHONY: infra/mailpit/shell
infra/mailpit/shell:
	$(PLATFORM_COMPOSE) exec $(MAILPIT_SERVICE) sh

.PHONY: infra/mailpit/remove
infra/mailpit/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(MAILPIT_SERVICE)
