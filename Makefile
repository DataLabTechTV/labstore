BIN_DIR := bin
BACKEND_DIR := backend
FRONTEND_DIR := web
BACKEND_CMD := $(BIN_DIR)/backend
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/dist
BENCHMARK_DIR := benchmark

.PHONY: all backend frontend build run benchmark clean

all: build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

BACKEND_SRCS := $(shell find $(BACKEND_DIR) -name "*.go")

$(BACKEND_CMD): $(BACKEND_SRCS) | $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(BACKEND_CMD) main.go

backend: $(BACKEND_CMD)

frontend:
	# cd $(FRONTEND_DIR) && npm install
	# cd $(FRONTEND_DIR) && npm run build

build: backend frontend

run: build
	set -a; . $(BACKEND_DIR)/.env; set +a; \
	(cd $(BACKEND_DIR) && ../$(BACKEND_CMD) serve --debug)# && \
	# (cd $(FRONTEND_DIR) && npm start)

benchmark:
	set -a; . $(BENCHMARK_DIR)/.env; \
	(cd $(BENCHMARK_DIR) && warp run config.yml)

clean:
	rm -rf $(BIN_DIR) $(FRONTEND_DIR)/node_modules $(FRONTEND_BUILD_DIR)
