BIN_DIR := bin

SERVER_DIR := server
CLIENT_DIR := client
ASSETS_DIR := $(SERVER_DIR)/pkg/router/assets

CLI_DIR := cli
CLI_CMD := $(BIN_DIR)/labstore

WEB_DIR := web
WEB_SRC_DIRS := $(WEB_DIR)/src $(WEB_DIR)/static
WEB_BUILD_DIR := $(WEB_DIR)/build

.PHONY: all web assets cli build run profile test clean-debug clean-cli clean-web clean

all: build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

WEB_SRCS := $(shell find $(WEB_SRC_DIRS) -type f)

$(WEB_BUILD_DIR): $(WEB_SRCS)
	cd $(WEB_DIR) && npm ci
	cd $(WEB_DIR) && npm run build

ASSETS_SRCS := $(shell find $(WEB_BUILD_DIR) -type f)

$(ASSETS_DIR): $(ASSETS_SRCS)
	rsync -a --delete web/build/ $(ASSETS_DIR)/

CLI_SRCS = $(shell find $(CLI_DIR) $(SERVER_DIR) $(CLIENT_DIR) -name '*.go')

$(CLI_CMD): $(CLI_SRCS) | $(BIN_DIR)
	cd $(CLI_DIR) && go build -o ../$(CLI_CMD) ./cmd/labstore

assets: web $(ASSETS_DIR)

cli: assets $(CLI_CMD)

web: $(WEB_BUILD_DIR)

build: cli

run: build
	npx concurrently \
		-n server,web \
		-c blue,green \
		"$(CLI_CMD) serve" \
		"cd $(WEB_DIR) && npm run preview -- --port 5123"

profile: cli
	npx concurrently \
		-n server,pprof \
		-c blue,yellow \
		"$(CLI_CMD) --pprof serve" \
		"go tool pprof \
			-focus=github.com/IllumiKnowLabs/labstore/server \
			-http=:8081 \
			http://localhost:6060/debug/pprof/profile?seconds=60"

SERVER_TEST_SRCS := $(shell find $(SERVER_DIR) -name '*.go')

test: $(SERVER_TEST_SRCS)
	cd $(SERVER_DIR) && go test ./... | grep -v '\[no test files\]'

clean-assets:
	rm -rf $(ASSETS_DIR)/

clean-debug:
	find . -type f -name '__debug_bin*' -delete

clean-cli: clean-assets clean-debug
	rm -rf $(BIN_DIR)/

clean-web:
	rm -rf $(WEB_DIR)/node_modules $(WEB_DIR)/.svelte-kit $(WEB_BUILD_DIR)

clean: clean-cli clean-web
