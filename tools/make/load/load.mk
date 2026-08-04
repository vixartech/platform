YANDEX_TANK_WORKDIR ?= $(CURDIR)/tools/load
YANDEX_TANK_CONFIG  ?= config.yaml
YANDEX_TANK_NETWORK ?= platform_default

.PHONY: load
load: load/seed
	@docker run \
		--rm \
		--platform linux/amd64 \
		--network $(YANDEX_TANK_NETWORK) \
		--volume "$(YANDEX_TANK_WORKDIR):/var/loadtest" \
		--workdir /var/loadtest \
		yandex/yandex-tank:latest \
		-c $(YANDEX_TANK_CONFIG)

.PHONY: load/seed
load/seed:
	@$(PLATFORM_COMPOSE) exec -T $(POSTGRES_SERVICE) \
		psql "$(POSTGRES_DSN)" -v ON_ERROR_STOP=1 < "$(YANDEX_TANK_WORKDIR)/seed.sql"
