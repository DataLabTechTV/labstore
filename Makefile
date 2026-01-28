BIN_DIR := bin

LABSTORE_CMD := $(BIN_DIR)/labstore

CMD_DIR := cmd
CLI_DIR := cli
CLIENT_DIR := client
SERVER_DIR := server
ASSETS_DIR := $(SERVER_DIR)/pkg/router/assets

WEB_DIR := web
WEB_SRC_DIRS := $(WEB_DIR)/src $(WEB_DIR)/static
WEB_BUILD_DIR := $(WEB_DIR)/build

HOST_GOOS := $(shell go env GOOS)
HOST_GOARCH := $(shell go env GOARCH)
HOST_GOARM := $(shell go env GOARM)

GOOS ?= $(HOST_GOOS)
GOARCH ?= $(HOST_GOARCH)
GOARM ?= $(HOST_GOARM)

ifeq ($(GOOS),windows)
    LABSTORE_CMD := $(LABSTORE_CMD).exe
endif

BIN_SUFFIX :=
ifeq ($(GOOS),$(HOST_GOOS))
  ifeq ($(GOARCH),$(HOST_GOARCH))
    BIN_SUFFIX :=
  else
    BIN_SUFFIX := -$(GOARCH)
    ifeq ($(GOARCH),arm)
      BIN_SUFFIX := $(BIN_SUFFIX)-$(GOARM)
    endif
  endif
else
  BIN_SUFFIX := -$(GOOS)-$(GOARCH)
  ifeq ($(GOARCH),arm)
    BIN_SUFFIX := $(BIN_SUFFIX)-$(GOARM)
  endif
endif

.PHONY: all
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

LABSTORE_SRCS = $(shell find $(CMD_DIR) $(CLI_DIR) $(SERVER_DIR) $(CLIENT_DIR) -name '*.go')

$(LABSTORE_CMD): $(LABSTORE_SRCS) | $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) \
	go build -o $(LABSTORE_CMD)$(BIN_SUFFIX) ./cmd/labstore

.PHONY: assets
assets: web $(ASSETS_DIR)

.PHONY: cli
cli: assets $(LABSTORE_CMD)

.PHONY: web
web: $(WEB_BUILD_DIR)

.PHONY: build-linux-amd64
build-linux-amd64:
	$(MAKE) build GOOS=linux GOARCH=amd64 LABSTORE_CMD=$(LABSTORE_CMD)-linux-amd64

.PHONY: build-linux-arm64
build-linux-arm64:
	$(MAKE) build GOOS=linux GOARCH=arm64 LABSTORE_CMD=$(LABSTORE_CMD)-linux-arm64

.PHONY: build-linux-armv7
build-linux-armv7:
	$(MAKE) build GOOS=linux GOARCH=arm GOARM=7 LABSTORE_CMD=$(LABSTORE_CMD)-linux-armv7

.PHONY: build-darwin-amd64
build-darwin-amd64:
	$(MAKE) build GOOS=darwin GOARCH=amd64 LABSTORE_CMD=$(LABSTORE_CMD)-darwin-amd64

.PHONY: build-darwin-arm64
build-darwin-arm64:
	$(MAKE) build GOOS=darwin GOARCH=arm64 LABSTORE_CMD=$(LABSTORE_CMD)-darwin-arm64

.PHONY: build-windows-amd64
build-windows-amd64:
	$(MAKE) build GOOS=windows GOARCH=amd64 LABSTORE_CMD=$(LABSTORE_CMD)-windows-amd64.exe

.PHONY: build-windows-arm64
build-windows-arm64:
	$(MAKE) build GOOS=windows GOARCH=arm64 LABSTORE_CMD=$(LABSTORE_CMD)-windows-arm64.exe

.PHONY: cli
build: cli

.PHONY: build-all
build-all: build-linux-amd64 build-linux-arm64 build-linux-armv7 \
	build-darwin-amd64 build-darwin-arm64 \
	build-windows-amd64 build-windows-arm64

.PHONY: run
run: build
	npx concurrently \
		-n server,web \
		-c blue,green \
		"$(LABSTORE_CMD) serve" \
		"cd $(WEB_DIR) && npm run preview -- --port 5123"

.PHONY: profile
profile: cli
	npx concurrently \
		-n server,pprof \
		-c blue,yellow \
		"$(LABSTORE_CMD) --pprof serve" \
		"go tool pprof \
			-focus=github.com/IllumiKnowLabs/labstore/server \
			-http=:8081 \
			http://localhost:6060/debug/pprof/profile?seconds=60"

SERVER_TEST_SRCS := $(shell find $(SERVER_DIR) -name '*.go')

.PHONY: test
test: $(SERVER_TEST_SRCS)
	cd $(SERVER_DIR) && go test ./... | grep -v '\[no test files\]'

.PHONY: clean-assets
clean-assets:
	rm -rf $(ASSETS_DIR)/

.PHONY: clean-cli
clean-cli: clean-assets
	rm -rf $(BIN_DIR)/

.PHONY: clean-web
clean-web:
	rm -rf $(WEB_DIR)/node_modules $(WEB_DIR)/.svelte-kit $(WEB_BUILD_DIR)

.PHONY: clean
clean: clean-cli clean-web

.PHONY: dist-clean
dist-clean: clean
	find . -type f -name '__debug_bin*' -delete
	go clean -cache
