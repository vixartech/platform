FROM node:24.19.0-alpine3.23 AS dependencies

WORKDIR /app
RUN corepack enable && corepack install --global pnpm@10.28.2
COPY . .
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --prod --frozen-lockfile --ignore-scripts

FROM node:24.19.0-alpine3.23

ENV NODE_ENV=production
WORKDIR /app
RUN apk add --no-cache chromium ca-certificates tini=0.19.0-r3
COPY --from=dependencies --chown=node:node /app/node_modules ./node_modules
COPY --chown=node:node . .
USER node
ENTRYPOINT ["/sbin/tini", "-g", "--", "node"]
