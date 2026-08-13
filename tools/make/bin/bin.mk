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

bin/buf: $(BIN_PATH)/buf

$(BIN_PATH)/buf: | $(BIN_PATH)
	@$(GOINSTALL) github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)

bin/migrate: $(BIN_PATH)/migrate

$(BIN_PATH)/migrate: | $(BIN_PATH)
	@$(GOINSTALL) -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

bin/easyjson: $(BIN_PATH)/easyjson

$(BIN_PATH)/easyjson: | $(BIN_PATH)
	@$(GOINSTALL) github.com/mailru/easyjson/...@$(EASYJSON_VERSION)

bin/oapi-codegen: $(BIN_PATH)/oapi-codegen

$(BIN_PATH)/oapi-codegen: | $(BIN_PATH)
	@$(GOINSTALL) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

bin/golangci-lint: $(BIN_PATH)/golangci-lint

$(BIN_PATH)/golangci-lint: | $(BIN_PATH)
	@wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh | \
		sh -s -- -b $(BIN_PATH) $(GOLANGCI_LINT_VERSION)

bin/sqlc: $(BIN_PATH)/sqlc

$(BIN_PATH)/sqlc: | $(BIN_PATH)
	@$(GOINSTALL) github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)

bin/protoc-gen-go: $(BIN_PATH)/protoc-gen-go

$(BIN_PATH)/protoc-gen-go: | $(BIN_PATH)
	@$(GOINSTALL) google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

bin/protoc-gen-go-grpc: $(BIN_PATH)/protoc-gen-go-grpc

$(BIN_PATH)/protoc-gen-go-grpc: | $(BIN_PATH)
	@$(GOINSTALL) google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

bin/protoc-gen-go-vtproto: $(BIN_PATH)/protoc-gen-go-vtproto

$(BIN_PATH)/protoc-gen-go-vtproto: | $(BIN_PATH)
	@$(GOINSTALL) github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@$(PROTOC_GEN_GO_VTPROTO_VERSION)

.PHONY: bin/clean
bin/clean:
	@rm -rf $(BIN_PATH)/*
