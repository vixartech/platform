DOCKER ?= docker

GO_DOCKERFILE ?= $(PLATFORM_DIR)/images/go/Dockerfile
GO_IMAGE_TAG  ?= local
GO_IMAGE_TOOLS ?=

go/image/%:
	@$(DOCKER) build \
		--file $(GO_DOCKERFILE) \
		--build-arg APP=$* \
		--build-arg TOOLS="$(GO_IMAGE_TOOLS)" \
		--tag $*:$(GO_IMAGE_TAG) \
		.
