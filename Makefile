PLATFORM_SERVICES := postgres redis nats listmonk nginx grafana alloy loki
PLATFORM_DIR      := $(CURDIR)

include platform.mk
