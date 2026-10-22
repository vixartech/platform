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
FROM debian:13-slim@sha256:62bd423aeee912639ca8ee681c2b762cb3e25137e0a92b20be455898a9f034ad AS browser

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        ca-certificates \
        chromium=155.0.8059.79-1~deb13u1 \
        chromium-common=155.0.8059.79-1~deb13u1 \
        tini=0.19.0-3+b8 \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 1001 app \
    && useradd --uid 1001 --gid app --create-home --shell /usr/sbin/nologin app

WORKDIR /app
COPY --from=builder /out/ ./
COPY configs configs

ENV XDG_CONFIG_HOME=/tmp/chromium/config \
    XDG_CACHE_HOME=/tmp/chromium/cache \
    XDG_DATA_HOME=/tmp/chromium/data
USER app
EXPOSE 8080 50051
ENTRYPOINT ["/usr/bin/tini", "-g", "--", "./app"]

# Preserve the existing default image for ordinary Go services.
FROM runtime AS app
