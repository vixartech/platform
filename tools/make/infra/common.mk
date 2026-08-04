.PHONY: infra/config
infra/config:
	$(PLATFORM_COMPOSE) config

.PHONY: infra/up
infra/up:
	$(PLATFORM_COMPOSE) up -d $(PLATFORM_SERVICES)

.PHONY: infra/start
infra/start:
	$(PLATFORM_COMPOSE) start $(PLATFORM_SERVICES)

.PHONY: infra/stop
infra/stop:
	$(PLATFORM_COMPOSE) stop $(PLATFORM_SERVICES)

.PHONY: infra/restart
infra/restart:
	$(PLATFORM_COMPOSE) restart $(PLATFORM_SERVICES)

.PHONY: infra/ps
infra/ps:
	$(PLATFORM_COMPOSE) ps

.PHONY: infra/logs
infra/logs:
	$(PLATFORM_COMPOSE) logs -f $(PLATFORM_SERVICES)

.PHONY: infra/down
infra/down:
	$(PLATFORM_COMPOSE) down --remove-orphans

.PHONY: infra/destroy
infra/destroy:
	$(PLATFORM_COMPOSE) down --volumes --remove-orphans
