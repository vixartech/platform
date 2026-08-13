PNPM ?= pnpm
WEB_DEV_PORT ?=
WEB_DEV_FLAGS ?=
WEB_CHECK_SCRIPT ?= check
WEB_BUILD_SCRIPT ?= build
WEB_E2E_SCRIPT ?= test:e2e

ifneq ($(strip $(WEB_DEV_PORT)),)
WEB_DEV_FLAGS += -p $(WEB_DEV_PORT)
endif

.PHONY: web/install
web/install:
	@$(PNPM) install

.PHONY: web/dev
web/dev:
	@$(PNPM) exec next dev $(WEB_DEV_FLAGS)

.PHONY: web/check
web/check:
	@$(PNPM) run $(WEB_CHECK_SCRIPT)

.PHONY: web/build
web/build:
	@$(PNPM) run $(WEB_BUILD_SCRIPT)

.PHONY: web/test/e2e
web/test/e2e:
	@$(PNPM) run $(WEB_E2E_SCRIPT)
