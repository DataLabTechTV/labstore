set shell := ["bash", "-uc"]

config_path := "labstore.yml"

[group("projects")]
mod backend "backend/justfile"

[group("projects")]
mod infra "infra/justfile"

[group("projects")]
mod benchmark "benchmark/justfile"

default:
    just -l

[group("helpers")]
check binary:
    @echo -n "Checking {{binary}}... "
    @command -v {{binary}} >/dev/null 2>&1 \
        || (echo "failed ({{binary}} not found)"; exit 1)
    @echo ok

[group("helpers")]
check-container container:
    @just check docker
    @echo -n "Checking running container {{container}}... "
    @docker ps --format '{{{{.Names}}' | grep -q '^{{container}}$' \
        || (echo "failed ({{container}} not running)"; exit 1)
    @echo ok

[group("helpers")]
check-port host port:
    @just check nc
    @echo -n "Checking for open port {{host}}:{{port}}... "
    @nc -z {{host}} {{port}} >/dev/null 2>&1 \
        || (echo "failed (closed)"; exit 1)
    @echo ok

[group("helpers")]
check-repo-deps:
    @just check pre-commit
    @just check golangci-lint

[group("helpers")]
check-deps:
    @just check-repo-deps
    @just backend check-deps
    @just infra check-deps
    @just benchmark check-deps

[group("helpers")]
print-env:
    #!/bin/bash

    LABSTORE_WEB_HOST=$(yq ".web.host" {{config_path}})
    LABSTORE_WEB_PORT=$(yq ".web.port" {{config_path}})

    LABSTORE_BACKEND_ADMIN_SERVER_HOST=$(yq ".backend.admin.server.host" {{config_path}})
    LABSTORE_BACKEND_ADMIN_SERVER_PORT=$(yq ".backend.admin.server.port" {{config_path}})
    LABSTORE_BACKEND_ADMIN_AUTH_ACCESS_KEY=$(yq ".backend.admin.auth.access_key" {{config_path}})
    LABSTORE_BACKEND_ADMIN_AUTH_SECRET_KEY=$(yq ".backend.admin.auth.secret_key" {{config_path}})

    LABSTORE_BACKEND_IAM_SERVER_HOST=$(yq ".backend.iam.server.host" {{config_path}})
    LABSTORE_BACKEND_IAM_SERVER_PORT=$(yq ".backend.iam.server.port" {{config_path}})

    LABSTORE_BACKEND_S3_SERVER_HOST=$(yq ".backend.s3.server.host" {{config_path}})
    LABSTORE_BACKEND_S3_SERVER_PORT=$(yq ".backend.s3.server.port" {{config_path}})

    BENCHMARK_HOST=$(yq ".benchmark.host" {{config_path}})

    BENCHMARK_PORTS_IPERF3=$(yq ".benchmark.ports.iperf3" {{config_path}})
    BENCHMARK_PORTS_LABSTORE=$(yq ".benchmark.ports.labstore" {{config_path}})
    BENCHMARK_PORTS_MINIO=$(yq ".benchmark.ports.minio" {{config_path}})
    BENCHMARK_PORTS_GARAGE=$(yq ".benchmark.ports.garage" {{config_path}})
    BENCHMARK_PORTS_SEAWEEDFS=$(yq ".benchmark.ports.seaweedfs" {{config_path}})
    BENCHMARK_PORTS_RUSTFS=$(yq ".benchmark.ports.rustfs" {{config_path}})

    BENCHMARK_STORE_ACCESS_KEY=$(yq ".benchmark.store.access_key" {{config_path}})
    BENCHMARK_STORE_SECRET_KEY=$(yq ".benchmark.store.secret_key" {{config_path}})
    BENCHMARK_STORE_TLS=$(yq ".benchmark.store.tls" {{config_path}})
    BENCHMARK_STORE_REGION=$(yq ".benchmark.store.region" {{config_path}})

    cat <<EOF
    export LABSTORE_WEB_HOST=$LABSTORE_WEB_HOST
    export LABSTORE_WEB_PORT=$LABSTORE_WEB_PORT
    LABSTORE_BACKEND_ADMIN_SERVER_HOST=$LABSTORE_BACKEND_ADMIN_SERVER_HOST
    LABSTORE_BACKEND_ADMIN_SERVER_PORT=$LABSTORE_BACKEND_ADMIN_SERVER_PORT
    LABSTORE_BACKEND_ADMIN_AUTH_ACCESS_KEY=$LABSTORE_BACKEND_ADMIN_AUTH_ACCESS_KEY
    LABSTORE_BACKEND_ADMIN_AUTH_SECRET_KEY=$LABSTORE_BACKEND_ADMIN_AUTH_SECRET_KEY
    LABSTORE_BACKEND_IAM_SERVER_HOST=$LABSTORE_BACKEND_IAM_SERVER_HOST
    LABSTORE_BACKEND_IAM_SERVER_PORT=$LABSTORE_BACKEND_IAM_SERVER_PORT
    LABSTORE_BACKEND_S3_SERVER_HOST=$LABSTORE_BACKEND_S3_SERVER_HOST
    LABSTORE_BACKEND_S3_SERVER_PORT=$LABSTORE_BACKEND_S3_SERVER_PORT
    export BENCHMARK_HOST=$BENCHMARK_HOST
    export BENCHMARK_PORTS_IPERF3=$BENCHMARK_PORTS_IPERF3
    export BENCHMARK_PORTS_LABSTORE=$BENCHMARK_PORTS_LABSTORE
    export BENCHMARK_PORTS_MINIO=$BENCHMARK_PORTS_MINIO
    export BENCHMARK_PORTS_GARAGE=$BENCHMARK_PORTS_GARAGE
    export BENCHMARK_PORTS_SEAWEEDFS=$BENCHMARK_PORTS_SEAWEEDFS
    export BENCHMARK_PORTS_RUSTFS=$BENCHMARK_PORTS_RUSTFS
    export BENCHMARK_STORE_ACCESS_KEY=$BENCHMARK_STORE_ACCESS_KEY
    export BENCHMARK_STORE_SECRET_KEY=$BENCHMARK_STORE_SECRET_KEY
    export BENCHMARK_STORE_TLS=$BENCHMARK_STORE_TLS
    export BENCHMARK_STORE_REGION=$BENCHMARK_STORE_REGION
    EOF

lint: check-repo-deps
    pre-commit run --all-files

install-hooks: check-repo-deps
    pre-commit install

clean:
    just benchmark clean
