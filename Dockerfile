# === Build: Go backend ======================================================
FROM golang:1.25-bookworm AS build-backend
WORKDIR /src
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/wiki ./cmd

# === Build: Nuxt frontend ===================================================
FROM node:22-bookworm-slim AS build-frontend
WORKDIR /app

# better-sqlite3 (via @nuxt/content) compiles native bindings at install time.
RUN apt-get update && apt-get install -y --no-install-recommends python3 make g++ \
    && rm -rf /var/lib/apt/lists/*

COPY docs-app/package.json docs-app/package-lock.json ./
RUN npm ci

COPY docs-app/ .
RUN npx nuxt build

# === Runtime ================================================================
# Node is required for the Nuxt server; the Go binary is statically linked.
FROM node:22-bookworm-slim
WORKDIR /app

# TLS roots so the JWT verifier can reach AOAuth's JWKS over HTTPS if needed.
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Go backend binary + seed data.
COPY --from=build-backend /out/wiki /app/wiki
COPY server/cmd/migration.yml /app/cmd/migration.yml
# godotenv.Load() fatals if .env is absent; bake a copy as a fallback.
# Runtime -e / environment: overrides take precedence over file values.

# Nuxt frontend build output (self-contained, includes compiled native bindings).
COPY --from=build-frontend /app/.output /app/.output

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Bind backend to all interfaces so Docker port-mapping reaches it.
# BACKEND_URL stays on loopback since both services share the container.
ENV HOST=0.0.0.0:1000 \
    NODE_ENV=production \
    NITRO_HOST=0.0.0.0 \
    NITRO_PORT=3000 \
    BACKEND_URL=http://localhost:1000

EXPOSE 3000 1000

ENTRYPOINT ["/app/entrypoint.sh"]
