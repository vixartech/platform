GOLANGCI_LINT_CONFIG ?= $(PLATFORM_DIR)/configs/go/golangci.yml
OAPI_CODEGEN_CONFIGS ?= $(wildcard api/http/*/*/config.yaml)
GOFMT_PATHS          ?= cmd internal pkg tools
GOARCH_FILES         := $(filter-out %_test.go,$(wildcard $(PLATFORM_DIR)/tools/goarch/*.go))

.PHONY: go/vendor
go/vendor:
	@go mod tidy
	@go mod vendor

.PHONY: go/generate
go/generate: bin/easyjson bin/mockgen
	@PATH="$(BIN_PATH):$$PATH" go generate ./...

.PHONY: go/generate/http
go/generate/http: bin/oapi-codegen bin/easyjson
	@for config in $(OAPI_CODEGEN_CONFIGS); do \
		$(BIN_PATH)/oapi-codegen -config "$$config" "$$(dirname "$$config")/api.yaml"; \
	done
	@for file in $$(find internal/pkg/http/gen -name gen.go -type f | sort); do \
		$(BIN_PATH)/easyjson -all -omit_empty "$$file"; \
	done

.PHONY: go/generate/proto
go/generate/proto: bin/buf bin/protoc-gen-go bin/protoc-gen-go-grpc bin/protoc-gen-go-vtproto
	@$(BIN_PATH)/buf dep update
	@$(BIN_PATH)/buf generate

.PHONY: go/generate/sqlc
go/generate/sqlc: bin/sqlc
	@$(BIN_PATH)/sqlc generate

.PHONY: go/architecture
go/architecture:
	@go run $(GOARCH_FILES)

.PHONY: go/lint
go/lint: bin/golangci-lint go/architecture
	@files="$$(for path in $(GOFMT_PATHS); do \
		[ ! -e "$$path" ] || find "$$path" -name '*.go' -type f; \
	done | sort)"; \
	if [ -n "$$files" ]; then gofmt -w $$files; fi
	@go vet ./...
	@$(BIN_PATH)/golangci-lint run ./... --config $(GOLANGCI_LINT_CONFIG) --fix

.PHONY: go/lint/proto
go/lint/proto: bin/buf
	@$(BIN_PATH)/buf lint

.PHONY: go/lint/proto/breaking
go/lint/proto/breaking: bin/buf
	@$(BIN_PATH)/buf breaking --against '.git#branch=master'

.PHONY: go/test
go/test:
	@go test ./...

.PHONY: go/test/race
go/test/race:
	@go test -race ./...

.PHONY: go/test/cover
go/test/cover:
	@profile="$$(mktemp "$${TMPDIR:-/tmp}/vixar-go-cover.XXXXXX")"; \
	trap 'rm -f "$$profile"' EXIT; \
	go test -coverprofile="$$profile" ./...; \
	go tool cover -func="$$profile"

.PHONY: go/vuln
go/vuln: bin/govulncheck
	@$(BIN_PATH)/govulncheck ./...
