DOCKER ?= docker

GO_DOCKERFILE ?= $(PLATFORM_DIR)/images/go/Dockerfile
GO_IMAGE_TAG  ?= local

go/image/%:
	@$(DOCKER) build \
		--file $(GO_DOCKERFILE) \
		--build-arg APP=$* \
		--tag $*:$(GO_IMAGE_TAG) \
		.
