FROM --platform=$BUILDPLATFORM golang:1.26.10-alpine3.24 AS builder

ARG APP
ARG TOOLS
ARG TARGETARCH
ARG TARGETOS

WORKDIR /app

COPY . .

# build app executable
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOARCH=$TARGETARCH GOOS=$TARGETOS go build \
    -mod=vendor \
    -trimpath \
    -ldflags="-s -w -buildid=" \
    -o /out/app \
    ./cmd/$APP && \
    for tool in $TOOLS; do \
      CGO_ENABLED=0 GOARCH=$TARGETARCH GOOS=$TARGETOS go build \
        -mod=vendor -trimpath -ldflags="-s -w -buildid=" \
        -o "/out/$tool" "./cmd/$tool" || exit 1; \
    done

FROM alpine:3.24 AS runtime

ARG APP

# create non-root user
RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app -u 1001 app

WORKDIR /app

COPY --from=builder /out/ ./
COPY configs configs

USER app

EXPOSE 8080 50051

ENTRYPOINT ["./app"]

# The explicit browser target is used only by services that render web pages.
FROM runtime AS browser

USER root
RUN apk add --no-cache chromium=155.0.8059.84-r0 tini=0.19.0-r3
ENV XDG_CONFIG_HOME=/tmp/chromium/config \
    XDG_CACHE_HOME=/tmp/chromium/cache
USER app
ENTRYPOINT ["/sbin/tini", "-g", "--", "./app"]

# Preserve the existing default image for ordinary Go services.
FROM runtime AS app
