FROM golang:1.26.5-alpine3.24 AS builder

ARG APP

WORKDIR /app

RUN apk add --no-cache bash make

COPY go.mod go.sum ./

# download app dependencies
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# build app executable
RUN --mount=type=cache,target=/root/.cache/go-build \
    make go/build/$APP GO_FLAGS=-mod=mod

FROM alpine:3.24

ARG APP

# create non-root user
RUN addgroup -S app && adduser -S -G app -u 1001 app

WORKDIR /app

COPY --from=builder /app/platform/bin/$APP app
COPY configs configs

USER app

EXPOSE 8080 50051

ENTRYPOINT ["./app"]
