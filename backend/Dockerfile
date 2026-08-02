# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy only the module files first so dependency download is cached
# and is not invalidated by unrelated source changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source
COPY ./ .

# Build the Go application
RUN go install -v ./...

EXPOSE 1323

# Command to run the executable
CMD ["api-oa-integrator"]
