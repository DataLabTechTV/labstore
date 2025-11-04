BIN_DIR := bin
BACKEND_DIR := backend
FRONTEND_DIR := web
BACKEND_CMD := $(BIN_DIR)/labstore-server
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/dist
BENCHMARK_DIR := benchmark
BENCHMARK_BUCKET := warp-benchmark-bucket

.PHONY: all backend frontend build run benchmark-bucket benchmark clean

all: build

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

BACKEND_SRCS := $(shell find $(BACKEND_DIR) -name "*.go")

$(BACKEND_CMD): $(BACKEND_SRCS) | $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(BACKEND_CMD) ./cmd/labstore-server

backend: $(BACKEND_CMD)

frontend:
	# cd $(FRONTEND_DIR) && npm install
	# cd $(FRONTEND_DIR) && npm run build

build: backend frontend

run: build
	set -a; . $(BACKEND_DIR)/.env; set +a; \
	(cd $(BACKEND_DIR) && ../$(BACKEND_CMD) serve --debug)# && \
	# (cd $(FRONTEND_DIR) && npm start)

benchmark-bucket:
	set -a; . $(BACKEND_DIR)/.env; set +a; \
	export MC_HOST_local=http:\/\/$${LS_ADMIN_ACCESS_KEY}:$${LS_ADMIN_SECRET_KEY}@$${LS_HOST}:$${LS_PORT}; \
	mc ls local/$(BENCHMARK_BUCKET) 2>&1 >/dev/null || mc mb local/$(BENCHMARK_BUCKET)

benchmark: benchmark-bucket
	set -a; . $(BENCHMARK_DIR)/.env; set +a; \
	(cd $(BENCHMARK_DIR) && mkdir -p output/ && cd output/ && warp run ../config.yml)

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(FRONTEND_DIR)/node_modules $(FRONTEND_BUILD_DIR)
	rm -rf $(BENCHMARK_DIR)/output
