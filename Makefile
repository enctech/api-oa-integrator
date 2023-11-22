BINARY_NAME=api-oa-integrator

# Set the default target (what happens when you run just `make`)
all: build

# The build target compiles your Go application
build:
	go build -o $(BINARY_NAME) main.go

# The clean target removes the compiled binary
clean:
	rm -f $(BINARY_NAME)

setup:
	tar -xzvf ssh.tar.gz