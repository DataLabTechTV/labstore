set shell := ["bash", "-uc"]

set dotenv-load
set dotenv-required

[group("main")]
mod backend "backend/justfile"

[group("main")]
mod infra "infra/justfile"

[group("main")]
mod benchmark "benchmark/justfile"

default:
    just -l

check binary:
    @echo -n "Checking {{binary}}... "
    @command -v {{binary}} >/dev/null 2>&1 \
        || (echo "failed ({{binary}} not found)"; exit 1)
    @echo ok

check-container container:
    @just check docker
    @echo -n "Checking running container {{container}}... "
    @docker ps --format '{{{{.Names}}' | grep -q '^{{container}}$' \
        || (echo "failed ({{container}} not running)"; exit 1)
    @echo ok

check-port host port:
    @just check nc
    @echo -n "Checking for open port {{host}}:{{port}}... "
    @nc -z {{host}} {{port}} >/dev/null 2>&1 \
        || (echo "failed ({{host}}:{{port}} closed)"; exit 1)
    @echo ok

clean:
    just benchmark clean
