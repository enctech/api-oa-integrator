# syntax=docker/dockerfile:1.4
FROM node:20 AS builder

# Set working directory
WORKDIR /usr/src/app

# Copy only package manager files first (improves caching)
COPY package.json pnpm-lock.yaml ./

# Install PNPM globally using npm
RUN npm install -g pnpm

# Install dependencies with cache
RUN --mount=type=cache,target=/root/.pnpm-store pnpm install --frozen-lockfile

# Copy the rest of the application files
COPY . .

# Build the app
RUN pnpm run build

# Use a lightweight production image
FROM node:20-alpine AS production

# Set working directory
WORKDIR /usr/src/app

# Copy built files from the builder stage
COPY --from=builder /usr/src/app .

# Expose the required port
EXPOSE 3000

# Run the app
CMD ["npx", "serve", "-s", "build"]
