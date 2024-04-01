BINARY_NAME=main
BUILD_DIR=./build
CMD_PATH=./cmd/server/main.go

.PHONY: build

build:
	@echo "Building $(BINARY_NAME) to $(BUILD_DIR)/$(BINARY_NAME)"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
