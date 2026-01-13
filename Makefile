BIN_DIR := bin

BACKEND_DIR := backend
BACKEND_CMD := $(BIN_DIR)/labstore-server

WEB_DIR := web
WEB_SRC_DIRS := $(WEB_DIR)/src $(WEB_DIR)/static
WEB_BUILD_DIR := $(WEB_DIR)/build

CLI_DIR := cli
CLI_CMD := $(BIN_DIR)/labstore

CLIENT_DIR := shared/client

.PHONY: all backend cli web build run profile test clean-debug clean-bin clean-web clean

all: build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

BACKEND_SRCS := $(shell find $(BACKEND_DIR) -name '*.go')

$(BACKEND_CMD): $(BACKEND_SRCS) | $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(BACKEND_CMD) ./cmd/labstore-server

backend: $(BACKEND_CMD)

WEB_SRCS := $(shell find $(WEB_SRC_DIRS) -type f)

$(WEB_BUILD_DIR): $(WEB_SRCS)
	cd $(WEB_DIR) && npm ci
	cd $(WEB_DIR) && npm run build

web: $(WEB_BUILD_DIR)

CLI_SRCS := $(shell find $(CLI_DIR) $(BACKEND_DIR) $(CLIENT_DIR) -name '*.go')

$(CLI_CMD): $(CLI_SRCS) | $(BIN_DIR)
	cd $(CLI_DIR) && go build -o ../$(CLI_CMD) ./cmd/labstore

cli: $(CLI_CMD)

build: backend cli web

run: build
	npx concurrently \
		-n backend,web \
		-c blue,green \
		"$(BACKEND_CMD) serve" \
		"cd $(WEB_DIR) && npm run preview -- --port 5123"

profile: backend
	npx concurrently \
		-n backend,pprof \
		-c blue,red \
		"$(BACKEND_CMD) --pprof serve" \
		"go tool pprof \
			-focus=github.com/IllumiKnowLabs/labstore/backend \
			-http=:8081 \
			http://localhost:6060/debug/pprof/profile?seconds=60"

test: $(BACKEND_TEST_SRCS)
	cd $(BACKEND_DIR) && go test ./... | grep -v '\[no test files\]'

clean-debug:
	find . -type f -name '__debug_bin*' -delete

clean-bin:
	rm -rf $(BIN_DIR)/

clean-web:
	rm -rf $(WEB_DIR)/node_modules $(WEB_DIR)/.svelte-kit $(WEB_BUILD_DIR)

clean: clean-debug clean-bin clean-web
