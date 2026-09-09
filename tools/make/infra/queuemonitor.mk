# Isolated, opt-in test stack. Not part of PLATFORM_COMPOSE or infra/up.
QUEUE_MONITOR_COMPOSE := $(COMPOSE) --project-name vixar-queue-monitor-local \
	-f $(PLATFORM_DIR)/deployments/queuemonitor/compose.yaml
export QUEUE_MONITOR_CONFIG_DIR

.PHONY: infra/queuemonitor/config infra/queuemonitor/up infra/queuemonitor/remove infra/queuemonitor/check
infra/queuemonitor/config:
	$(QUEUE_MONITOR_COMPOSE) config --quiet

infra/queuemonitor/up: infra/queuemonitor/config
	$(QUEUE_MONITOR_COMPOSE) up -d --wait --wait-timeout 60

# Removes only disposable state of this dedicated project, never the shared platform stack.
infra/queuemonitor/remove:
	$(QUEUE_MONITOR_COMPOSE) down --volumes

infra/queuemonitor/check: infra/queuemonitor/config
	$(QUEUE_MONITOR_COMPOSE) run --rm --no-deps --entrypoint promtool prometheus check config /config/prometheus.yaml
	$(QUEUE_MONITOR_COMPOSE) run --rm --no-deps --entrypoint promtool prometheus test rules /config/alerts.test.yaml
