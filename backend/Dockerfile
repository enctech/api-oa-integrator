# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY ./ .

# Build the Go application
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 1323

# Command to run the executable
CMD ["api-oa-integrator"]
