BIN_PATH ?= $(PLATFORM_DIR)/bin

GO_FLAGS ?= -mod=vendor
GO_BUILD_FLAGS ?= -trimpath
GO_LDFLAGS ?= -s -w -buildid=

go/build/%: | $(BIN_PATH)
	CGO_ENABLED=0 GOFLAGS="$(GO_FLAGS)" go build $(GO_BUILD_FLAGS) -ldflags="$(GO_LDFLAGS)" -o $(BIN_PATH)/$* ./cmd/$*

# platform/bin/gost --config configs/gost.yaml
