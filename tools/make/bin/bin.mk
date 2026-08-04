BIN_PATH ?= $(PLATFORM_DIR)/bin

BUF_VERSION                   ?= v1.72.0
MIGRATE_VERSION               ?= v4.19.1
EASYJSON_VERSION              ?= v0.9.2
OAPI_CODEGEN_VERSION          ?= v2.8.0
GOLANGCI_LINT_VERSION         ?= v2.10.1
SQLC_VERSION                  ?= v1.30.0
PROTOC_GEN_GO_VERSION         ?= v1.36.11
PROTOC_GEN_GO_GRPC_VERSION    ?= v1.6.2
PROTOC_GEN_GO_VTPROTO_VERSION ?= v0.6.0

GOINSTALL ?= GOBIN="$(BIN_PATH)" go install

$(BIN_PATH):
	@mkdir -p $(BIN_PATH)

.PHONY: bin/buf
bin/buf: | $(BIN_PATH)
	@$(GOINSTALL) github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)

.PHONY: bin/migrate
bin/migrate: | $(BIN_PATH)
	@$(GOINSTALL) -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

.PHONY: bin/easyjson
bin/easyjson: | $(BIN_PATH)
	@$(GOINSTALL) github.com/mailru/easyjson/...@$(EASYJSON_VERSION)

.PHONY: bin/oapi-codegen
bin/oapi-codegen: | $(BIN_PATH)
	@$(GOINSTALL) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

.PHONY: bin/golangci-lint
bin/golangci-lint: | $(BIN_PATH)
	@wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh | \
		sh -s -- -b $(BIN_PATH) $(GOLANGCI_LINT_VERSION)

.PHONY: bin/sqlc
bin/sqlc: | $(BIN_PATH)
	@$(GOINSTALL) github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)

.PHONY: bin/protoc-gen-go
bin/protoc-gen-go: | $(BIN_PATH)
	@$(GOINSTALL) google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

.PHONY: bin/protoc-gen-go-grpc
bin/protoc-gen-go-grpc: | $(BIN_PATH)
	@$(GOINSTALL) google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

.PHONY: bin/protoc-gen-go-vtproto
bin/protoc-gen-go-vtproto: | $(BIN_PATH)
	@$(GOINSTALL) github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@$(PROTOC_GEN_GO_VTPROTO_VERSION)

.PHONY: bin/clean
bin/clean:
	@rm -rf $(BIN_PATH)/*
