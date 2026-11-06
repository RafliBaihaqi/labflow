# syntax=docker/dockerfile:1.6

# --- Dependencies (cached) ---
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

# --- Build ---
FROM node:20-alpine AS build-env
WORKDIR /app
# Build-time public env - must be set for prod, fallback to Cloudflare test key for dev/CI
ARG VITE_TURNSTILE_SITE_KEY=1x00000000000000000000AA
ENV VITE_TURNSTILE_SITE_KEY=$VITE_TURNSTILE_SITE_KEY
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN echo "Building with VITE_TURNSTILE_SITE_KEY=${VITE_TURNSTILE_SITE_KEY:0:10}..." && npm run build

# --- Production dependencies ---
FROM node:20-alpine AS production-deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# --- Runner ---
FROM node:20-alpine AS runner

LABEL org.opencontainers.image.title="labflow-frontend"
LABEL org.opencontainers.image.description="LabFlow Go Frontend - React Router SSR"
LABEL org.opencontainers.image.source="https://github.com/RafliBaihaqi/labflow"

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000

# Install curl + wget for healthcheck + dumb-init for signal handling
RUN apk --no-cache add wget curl dumb-init

# Copy production node_modules and build artifacts
COPY --from=production-deps /app/node_modules ./node_modules
COPY --from=build-env /app/build ./build
COPY --from=build-env /app/package.json ./package.json

# Use non-root user (node user exists in node:alpine)
RUN chown -R node:node /app
USER node

EXPOSE 3000

HEALTHCHECK --interval=15s --timeout=10s --start-period=30s --retries=5 \
  CMD curl -f http://localhost:3000/ || wget -qO- http://localhost:3000/ > /dev/null || exit 1

ENTRYPOINT ["dumb-init", "--"]
CMD ["npm", "run", "start"]
