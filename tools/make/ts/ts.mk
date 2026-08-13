TS_DEV_PORT  ?= 3001
TS_DEV_FLAGS ?=

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

.PHONY: ts/build
ts/build:
	@pnpm run build

.PHONY: ts/test/e2e
ts/test/e2e:
	@pnpm run test:e2e
