NGINX_SERVICE ?= nginx

NGINX_IMAGE  ?= 1.30.4-alpine
NGINX_PORT   ?= 80
NGINX_CONFIG ?= $(PLATFORM_DIR)/deployments/nginx/configs/nginx.conf

override NGINX_CONFIG := $(abspath $(NGINX_CONFIG))

export NGINX_IMAGE
export NGINX_PORT
export NGINX_CONFIG

NGINX_COMPOSE_FILE := $(PLATFORM_DIR)/deployments/nginx/compose.yaml
PLATFORM_COMPOSE   += -f $(NGINX_COMPOSE_FILE)

.PHONY: infra/nginx/config
infra/nginx/config:
	$(PLATFORM_COMPOSE) config $(NGINX_SERVICE)

.PHONY: infra/nginx/up
infra/nginx/up:
	$(PLATFORM_COMPOSE) up -d $(NGINX_SERVICE)

.PHONY: infra/nginx/start
infra/nginx/start:
	$(PLATFORM_COMPOSE) start $(NGINX_SERVICE)

.PHONY: infra/nginx/stop
infra/nginx/stop:
	$(PLATFORM_COMPOSE) stop $(NGINX_SERVICE)

.PHONY: infra/nginx/restart
infra/nginx/restart:
	$(PLATFORM_COMPOSE) restart $(NGINX_SERVICE)

.PHONY: infra/nginx/ps
infra/nginx/ps:
	$(PLATFORM_COMPOSE) ps $(NGINX_SERVICE)

.PHONY: infra/nginx/logs
infra/nginx/logs:
	$(PLATFORM_COMPOSE) logs -f $(NGINX_SERVICE)

.PHONY: infra/nginx/shell
infra/nginx/shell:
	$(PLATFORM_COMPOSE) exec $(NGINX_SERVICE) sh

.PHONY: infra/nginx/test
infra/nginx/test:
	$(PLATFORM_COMPOSE) exec $(NGINX_SERVICE) nginx -t

.PHONY: infra/nginx/reload
infra/nginx/reload:
	$(PLATFORM_COMPOSE) exec $(NGINX_SERVICE) nginx -t
	$(PLATFORM_COMPOSE) exec $(NGINX_SERVICE) nginx -s reload

.PHONY: infra/nginx/remove
infra/nginx/remove:
	$(PLATFORM_COMPOSE) rm --stop --force $(NGINX_SERVICE)
