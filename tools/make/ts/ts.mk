TS_DEV_PORT  ?= 3001
TS_DEV_FLAGS ?=
TS_BIOME_CONFIG ?= $(PLATFORM_DIR)/configs/ts/biome.json

ifneq ($(strip $(TS_DEV_PORT)),)
TS_DEV_FLAGS += -p $(TS_DEV_PORT)
endif

.PHONY: ts/install
ts/install:
	@pnpm install

.PHONY: ts/dev
ts/dev:
	@pnpm exec next dev $(TS_DEV_FLAGS)

.PHONY: ts/check
ts/check:
	@pnpm run check:ci

.PHONY: ts/biome/check
ts/biome/check:
	@pnpm exec biome check . --config-path $(TS_BIOME_CONFIG)

.PHONY: ts/biome/write
ts/biome/write:
	@pnpm exec biome check . --write --config-path $(TS_BIOME_CONFIG)

.PHONY: ts/format
ts/format:
	@pnpm exec biome format . --write --config-path $(TS_BIOME_CONFIG)

.PHONY: ts/format/check
ts/format/check:
	@pnpm exec biome format . --config-path $(TS_BIOME_CONFIG)

.PHONY: ts/lint
ts/lint:
	@pnpm exec biome lint . --write --config-path $(TS_BIOME_CONFIG)

.PHONY: ts/lint/check
ts/lint/check:
	@pnpm exec biome lint . --config-path $(TS_BIOME_CONFIG)

.PHONY: ts/build
ts/build:
	@pnpm run build

.PHONY: ts/test/e2e
ts/test/e2e:
	@pnpm run test:e2e
