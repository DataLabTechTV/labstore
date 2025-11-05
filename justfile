set shell := ["bash", "-uc"]

set dotenv-load
set dotenv-required

mod backend "backend/justfile"
mod infra "infra/justfile"

default:
    just -l
