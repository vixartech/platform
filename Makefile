PLATFORM_SERVICES := postgres redis nats nginx grafana alloy loki
PLATFORM_DIR      := $(CURDIR)

include platform.mk
