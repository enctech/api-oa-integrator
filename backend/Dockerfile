# Stage 1: Build the Go application
FROM golang:latest AS builder

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
FROM alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /go/bin/api-oa-integrator /usr/local/bin

# Disable cache for this specific layer
RUN ls -la /usr/local/bin

EXPOSE 1323

# Command to run the executable
ENTRYPOINT ["/usr/local/bin/api-oa-integrator"]
