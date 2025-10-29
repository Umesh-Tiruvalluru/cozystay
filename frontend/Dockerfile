# ---- Stage 1: Build ----
FROM node:20-bookworm AS builder

# Set working directory
WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files first for caching
COPY pnpm-lock.yaml package.json ./

# Install dependencies (with frozen lockfile for reproducibility)
RUN pnpm install --frozen-lockfile

# Copy the rest of the project
COPY . .

# Build Next.js app
RUN pnpm build

# ---- Stage 2: Production ----
FROM node:20-bookworm-slim AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy necessary files from builder
COPY --from=builder /app/package.json ./
COPY --from=builder /app/pnpm-lock.yaml ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

# Expose port
EXPOSE 3000

# Run Next.js
CMD ["pnpm", "run", "dev"]
