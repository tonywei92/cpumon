# Variables
BINARY_NAME=cpumon

# Default target
all: build

# Build the Go project
build:
	@echo "Building the project..."
	go build -o $(BINARY_NAME) main.go

# Run the Go project
run:
	@echo "Running the project..."
	go run main.go

# Clean up the binary
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

.PHONY: all build run test clean logs
