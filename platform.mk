COMPOSE ?= docker compose

PLATFORM_SERVICES ?=
PLATFORM_DIR      ?= $(CURDIR)/platform

PLATFORM_COMPOSE := $(COMPOSE) \
	--project-name platform \
	--project-directory $(PLATFORM_DIR)

include $(PLATFORM_DIR)/tools/make/bin/bin.mk

include $(PLATFORM_DIR)/tools/make/infra/postgres.mk
include $(PLATFORM_DIR)/tools/make/infra/redis.mk
include $(PLATFORM_DIR)/tools/make/infra/nats.mk
include $(PLATFORM_DIR)/tools/make/infra/mailpit.mk
include $(PLATFORM_DIR)/tools/make/infra/listmonk.mk
include $(PLATFORM_DIR)/tools/make/infra/nginx.mk
include $(PLATFORM_DIR)/tools/make/infra/loki.mk
include $(PLATFORM_DIR)/tools/make/infra/alloy.mk
include $(PLATFORM_DIR)/tools/make/infra/grafana.mk
include $(PLATFORM_DIR)/tools/make/infra/queuemonitor.mk
include $(PLATFORM_DIR)/tools/make/infra/common.mk

include $(PLATFORM_DIR)/tools/make/go/build.mk
include $(PLATFORM_DIR)/tools/make/go/image.mk
include $(PLATFORM_DIR)/tools/make/go/dev.mk

include $(PLATFORM_DIR)/tools/make/ts/ts.mk

include $(PLATFORM_DIR)/tools/make/load/load.mk
