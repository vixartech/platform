TS_DEV_PORT  ?= 3001

.PHONY: ts/install
ts/install:
	@pnpm install

.PHONY: ts/dev
ts/dev:
	@pnpm exec next dev -p $(TS_DEV_PORT)

.PHONY: ts/check
ts/check:
	@pnpm run check:ci

.PHONY: ts/typecheck
ts/typecheck:
	@pnpm exec tsc --noEmit --incremental false

.PHONY: ts/lib/check
ts/lib/check: ts/biome/check ts/typecheck

.PHONY: ts/next/typegen
ts/next/typegen:
	@pnpm exec next typegen

.PHONY: ts/next/check
ts/next/check: ts/biome/check ts/next/typegen ts/typecheck

.PHONY: ts/biome/check
ts/biome/check:
	@pnpm exec biome check . --config-path $(PLATFORM_DIR)/configs/ts/biome.json

.PHONY: ts/biome/write
ts/biome/write:
	@pnpm exec biome check . --write --config-path $(PLATFORM_DIR)/configs/ts/biome.json

.PHONY: ts/format
ts/format:
	@pnpm exec biome format . --write --config-path $(PLATFORM_DIR)/configs/ts/biome.json

.PHONY: ts/format/check
ts/format/check:
	@pnpm exec biome format . --config-path $(PLATFORM_DIR)/configs/ts/biome.json

.PHONY: ts/lint
ts/lint:
	@pnpm exec biome lint . --write --config-path $(PLATFORM_DIR)/configs/ts/biome.json

.PHONY: ts/lint/check
ts/lint/check:
	@pnpm exec biome lint . --config-path $(PLATFORM_DIR)/configs/ts/biome.json

.PHONY: ts/build
ts/build:
	@pnpm run build

.PHONY: ts/next/build
ts/next/build:
	@pnpm exec next build

.PHONY: ts/next/start
ts/next/start:
	@pnpm exec next start

.PHONY: ts/test/e2e
ts/test/e2e:
	@pnpm run test:e2e
