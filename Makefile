VERSION := 0.1.0
BINARY_NAME := "rlss"
PACKAGE_NAME := "realess-server"
DISPLAY_NAME := "Realess Server"
BUILD_DIR := "./release/build"
RELEASE_DIR := "./release"

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))


build:
	@mkdir -p $(BUILD_DIR) && \
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/rlss/main.go

pack:
	@mkdir -p $(RELEASE_DIR) && \
	cp ./setup.sh $(BUILD_DIR)/ && \
	makeself $(BUILD_DIR) $(RELEASE_DIR)/$(PACKAGE_NAME)-Linux-amd64.sh $(DISPLAY_NAME) ./setup.sh

# This command is used for test like: fwg ...
run:
	@go run ./cmd/rlss/main.go $(ARGS)


# To prevent make from attempting to build a second target, add the catch-all rule
%:
	@:
