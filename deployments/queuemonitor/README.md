# Isolated queue-monitoring laboratory

Opt-in infrastructure only: disposable Redis, Redis exporter, Prometheus,
Alertmanager and a bounded local webhook sink. No application or crawler is run.
It is **not** part of `PLATFORM_COMPOSE` or ordinary `infra/up`.

From a consuming repository, set `PLATFORM_DIR` to this checkout and
`QUEUE_MONITOR_CONFIG_DIR` to an absolute product-owned directory containing:

- `queue.lua`: metadata collection, executed with `EVAL_RO`;
- `prometheus.yaml`, `alerts.yaml`, `alerts.test.yaml`: scraping and rule tests;
- `alertmanager.yaml`: routing exclusively to `http://receiver:8080/alerts`.

Use `infra/queuemonitor/config`, `infra/queuemonitor/check`, and
`infra/queuemonitor/up`. `check` runs native Prometheus validation and rule tests.
These targets require Docker Compose and pinned images; missing images are pulled
explicitly by Compose. Do not use untrusted product config or real credentials.

Project name is `vixar-queue-monitor-local`. Its network is internal; no host ports
or shared infrastructure networks are used. Redis has no authentication only
inside this disposable isolated network. Images are pinned by digest. Container
memory and tmpfs are bounded; Prometheus retention is 30 minutes/32 MiB. These
retention settings are not a hard limit on total TSDB overhead: tmpfs and container
memory limits provide the upper bound. Redis is limited to 16 MiB/noeviction.

Query Prometheus/exporter/receiver through `docker compose exec -T receiver node`
inside the network. The webhook holds only the last 100 alert names/statuses/times
in memory. Bodies are capped at 64 KiB, connections at eight; labels, annotations
and payloads are never stored. This is a test sink, not a delivery channel for people.

`infra/queuemonitor/remove` removes this project's containers/network/volumes,
discarding its synthetic Redis state, TSDB and notification history. It never
removes the ordinary `platform` project. No backup is made for this disposable lab.
Do not reuse its fixed project name for valuable or unrelated state.
