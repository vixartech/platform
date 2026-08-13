PNPM ?= pnpm
TS_DEV_PORT ?=
TS_DEV_FLAGS ?=
TS_CHECK_SCRIPT ?= check
TS_BUILD_SCRIPT ?= build
TS_E2E_SCRIPT ?= test:e2e

ifneq ($(strip $(TS_DEV_PORT)),)
TS_DEV_FLAGS += -p $(TS_DEV_PORT)
endif

.PHONY: ts/install
ts/install:
	@$(PNPM) install

.PHONY: ts/dev
ts/dev:
	@$(PNPM) exec next dev $(TS_DEV_FLAGS)

.PHONY: ts/check
ts/check:
	@$(PNPM) run $(TS_CHECK_SCRIPT)

.PHONY: ts/build
ts/build:
	@$(PNPM) run $(TS_BUILD_SCRIPT)

.PHONY: ts/test/e2e
ts/test/e2e:
	@$(PNPM) run $(TS_E2E_SCRIPT)
