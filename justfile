set shell := ["bash", "-uc"]

set dotenv-load
set dotenv-required

mod backend "backend/justfile"
mod infra "infra/justfile"
mod benchmark "benchmark/justfile"

default:
    just -l

clean:
    just benchmark clean
