# Stage 1: Build the Go application
FROM golang:1.23 AS builder

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy go.mod and go.sum files to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy the rest of the application source code
COPY ./ .

# Build the Go application with cached dependencies and build cache
RUN --mount=type=cache,target=/root/.cache/go-build go build -mod=readonly -o /go/bin/api-oa-integrator

# Stage 2: Create a smaller image for running the application
FROM debian:bookworm-slim

# Links the published package to this repository, which is what lets
# GITHUB_TOKEN publish it under the repo owner's namespace.
LABEL org.opencontainers.image.source=https://github.com/enctech/api-oa-integrator

# Copy the built binary from the builder stage
COPY --from=builder /go/bin/api-oa-integrator /usr/local/bin
COPY --from=builder /go/src/app/database/migrations /app/database/migrations

EXPOSE 1323

# ENTRYPOINT, not CMD: docker-compose passes --config and --migrations as
# arguments. With a bare CMD those would replace the binary instead.
ENTRYPOINT ["/usr/local/bin/api-oa-integrator"]
