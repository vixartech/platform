REDIS_SERVICE ?= redis

REDIS_IMAGE    ?= 8.8.1-alpine
REDIS_PASSWORD ?= password
REDIS_PORT     ?= 6379
REDIS_CONFIG   ?= $(PLATFORM_DIR)/deployments/redis/configs/redis.conf

export REDIS_IMAGE
export REDIS_PASSWORD
export REDIS_PORT
export REDIS_CONFIG

REDIS_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/redis/compose.yaml
PLATFORM_COMPOSE   += -f $(REDIS_COMPOSE_FILE)

.PHONY: infra/redis/config
infra/redis/config:
	$(PLATFORM_COMPOSE) config $(REDIS_SERVICE)

.PHONY: infra/redis/up
infra/redis/up:
	$(PLATFORM_COMPOSE) up -d $(REDIS_SERVICE)

.PHONY: infra/redis/start
infra/redis/start:
	$(PLATFORM_COMPOSE) start $(REDIS_SERVICE)

.PHONY: infra/redis/stop
infra/redis/stop:
	$(PLATFORM_COMPOSE) stop $(REDIS_SERVICE)

.PHONY: infra/redis/restart
infra/redis/restart:
	$(PLATFORM_COMPOSE) restart $(REDIS_SERVICE)

.PHONY: infra/redis/ps
infra/redis/ps:
	$(PLATFORM_COMPOSE) ps $(REDIS_SERVICE)

.PHONY: infra/redis/logs
infra/redis/logs:
	$(PLATFORM_COMPOSE) logs -f $(REDIS_SERVICE)

.PHONY: infra/redis/shell
infra/redis/shell:
	$(PLATFORM_COMPOSE) exec -e REDISCLI_AUTH="$(REDIS_PASSWORD)" $(REDIS_SERVICE) redis-cli

.PHONY: infra/redis/remove
infra/redis/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(REDIS_SERVICE)
