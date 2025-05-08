# Variables
BINARY_NAME=demoparser.exe
SRC_DIR=.

# Build the binary
build:
	go build -o bin/$(BINARY_NAME) $(SRC_DIR)

# Clean up the binary
clean:
	rm -f $(BINARY_NAME)

# Format the code
fmt:
	go fmt ./...

# Run tests
test:
	go test ./...

# Tidy up dependencies
tidy:
	go mod tidy

# Install dependencies
deps:
	go mod download

# # Help command to list all available commands
# help:
#     @echo "Available commands:"
#     @echo "  make build   - Build the binary"
#     @echo "  make run     - Build and run the program"
#     @echo "  make clean   - Remove the binary"
#     @echo "  make fmt     - Format the code"
#     @echo "  make test    - Run tests"
#     @echo "  make tidy    - Tidy up dependencies"
#     @echo "  make deps    - Install dependencies"