BIN_DIR := bin

BACKEND_DIR := backend
BACKEND_CMD := $(BIN_DIR)/labstore-server

FRONTEND_DIR := web
FRONTEND_SRC_DIRS := $(FRONTEND_DIR)/src $(FRONTEND_DIR)/static
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/build

CLI_DIR := cli
CLI_CMD := $(BIN_DIR)/labstore

CLIENT_DIR := shared/client

.PHONY: all backend frontend cli build run profile test clean

all: build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

BACKEND_SRCS := $(shell find $(BACKEND_DIR) -name '*.go')

$(BACKEND_CMD): $(BACKEND_SRCS) | $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(BACKEND_CMD) ./cmd/labstore-server

backend: $(BACKEND_CMD)

FRONTEND_SRCS := $(shell find $(FRONTEND_SRC_DIRS) -type f)

$(FRONTEND_BUILD_DIR): $(FRONTEND_SRCS)
	cd $(FRONTEND_DIR) && npm ci
	cd $(FRONTEND_DIR) && npm run build

frontend: $(FRONTEND_BUILD_DIR)

CLI_SRCS := $(shell find $(CLI_DIR) $(BACKEND_DIR) $(CLIENT_DIR) -name '*.go')

$(CLI_CMD): $(CLI_SRCS) | $(BIN_DIR)
	cd $(CLI_DIR) && go build -o ../$(CLI_CMD) ./cmd/labstore

cli: $(CLI_CMD)

build: backend frontend cli

run: build
	npx concurrently \
		-n backend,web \
		-c blue,green \
		"$(BACKEND_CMD) serve" \
		"cd $(FRONTEND_DIR) && npm run preview -- --port 5123"

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

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(FRONTEND_DIR)/node_modules $(FRONTEND_DIR)/.svelte-kit $(FRONTEND_BUILD_DIR)
