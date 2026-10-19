FROM --platform=$BUILDPLATFORM golang:1.26.9-alpine3.24 AS builder

ARG APP
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
    ./cmd/$APP

FROM alpine:3.24

ARG APP

# create non-root user
RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app -u 1001 app

WORKDIR /app

COPY --from=builder /out/app app
COPY configs configs

USER app

EXPOSE 8080 50051

ENTRYPOINT ["./app"]
