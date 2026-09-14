FROM node:24.19.0-alpine3.23 AS base

WORKDIR /app

RUN apk add --no-cache git openssh-client \
    && corepack enable \
    && corepack install --global pnpm@10.28.2 \
    && printf '%s\n' \
    'github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl' \
    > /etc/ssh/ssh_known_hosts

FROM base AS dependencies

# Workspace packages and their manifests must exist when pnpm resolves workspace links.
# Products own a restrictive .dockerignore so only build inputs enter this stage.
COPY . .

RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    --mount=type=ssh,required=true \
    git config --global url."ssh://git@github.com/".insteadOf "https://github.com/" \
    && GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=yes" \
    pnpm install --frozen-lockfile --ignore-scripts

FROM base AS builder

ARG API_URL
ARG ID_API_URL
ARG NEXT_PUBLIC_ID_API_BASE_PATH
ARG NEXT_PUBLIC_ID_OAUTH_BASE_URL
ARG NEXT_PUBLIC_ID_PROJECT_ID
ARG NEXT_PUBLIC_SHORT_LINK_URL

ENV API_URL=$API_URL \
    ID_API_URL=$ID_API_URL \
    NEXT_PUBLIC_ID_API_BASE_PATH=$NEXT_PUBLIC_ID_API_BASE_PATH \
    NEXT_PUBLIC_ID_OAUTH_BASE_URL=$NEXT_PUBLIC_ID_OAUTH_BASE_URL \
    NEXT_PUBLIC_ID_PROJECT_ID=$NEXT_PUBLIC_ID_PROJECT_ID \
    NEXT_PUBLIC_SHORT_LINK_URL=$NEXT_PUBLIC_SHORT_LINK_URL \
    NEXT_TELEMETRY_DISABLED=1

COPY --from=dependencies /app/node_modules ./node_modules
COPY . .

RUN mkdir -p public \
    && pnpm exec next build

FROM node:24.19.0-alpine3.23 AS runtime

ENV HOSTNAME=0.0.0.0 \
    NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000

WORKDIR /app

COPY --from=builder --chown=node:node /app/.next/standalone ./
COPY --from=builder --chown=node:node /app/.next/static ./.next/static
COPY --from=builder --chown=node:node /app/public ./public

USER node

EXPOSE 3000

CMD ["node", "server.js"]
