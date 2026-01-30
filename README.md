<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/IllumiKnowLabs/.github/refs/heads/main/branding/logos/SVG/projects/labstore/logo-dark.svg">
  <img alt="Shows the IllumiKnow Labs logo, featuring a black hole with a colorful accretion disk emitting light from the singularity. Above it, the acronym IKL." src="https://raw.githubusercontent.com/IllumiKnowLabs/.github/refs/heads/main/branding/logos/SVG/projects/labstore/logo-light.svg" width="150px">
</picture>


# LabStore

An object store with S3- and IAM-compatible APIs, built to support system prototyping and data-driven projects.


## Installation

We try to make installation as accessible as possible, as to not waste your time. Depending on your preferred approach, we provide:

- Docker images, available on [GitHub Container Registry (GHCR)](https://github.com/IllumiKnowLabs/labstore/pkgs/container/labstore/versions?filters%5Bversion_type%5D=tagged).
- A working `go install` command.
- Prebuilt binaries, available on [GitHub Releases](https://github.com/IllumiKnowLabs/labstore/releases).


### Docker

#### Prebuilt Images

To start LabStore server from our pre-built image:

```bash
docker run -d -p 6789:6789 ghcr.io/illumiknowlabs/labstore:v0.1.0 serve
```

We also recommend that you add an alias to your shell startup file:

```bash
alias labstore='docker run -it ghcr.io/illumiknowlabs/labstore:v0.1.0'
```

And then you can use the CLI, based on the the same image:

```bash
labstore s3 service list-buckets
```

#### Build Your Own

If you prefer, you can also build your own image and customize your own container using Docker Compose alongside available configuration environment variables and CLI arguments. See [Configuration](#configuration) for more details.

Just clone this repo and take a look at the `labstore` service under [infra/compose.yml](infra/compose.yml). In order to build the image, make sure you use the corresponding `just` command, since we the existing `compose.yml` depends on several environment variables that we load there:

```bash
just infra::up
```

This will launch the server and also provide the image you can use as a client:

```bash
alias labstore='docker run -it labstore-labstore'
```


### go install

If you already have Go installed, you might prefer to use that to install LabStore:

```bash
go install github.com/IllumiKnowLabs/labstore/cmd/labstore@v0.1.0
```

Please notice that this version does not come with the web UI embedded, but it will still work without an issue by downloading the required files from the corresponding GitHub Release during runtime.

You can check whether your `labstore` binary uses embedded or runtime web UI assets by running:

```bash
labstore version
```

You should see something like this (notice the last line):

```
╭────────────────────────────────────────────╮
│ 🚀 Welcome to LabStore, by IllumiKnow Labs │
╰────────────────────────────────────────────╯
version: 0.1.0 (tag: unknown, commit: unknown)
build time: unknown, builder: unknown
embedded web assets: no
```

> [!NOTE]
> This is the least preferred installation approach, since it won't be able to set `ldflags`, thus not tracking version information properly (as seen above), and it also does not embed the web UI assets, making it slower on first run and dependent on an internet connection being available on the first run.


### Prebuilt Binaries

Finally, we also provide prebuilt binaries that you can download directly from the desired version under [GitHub Releases](https://github.com/IllumiKnowLabs/labstore/releases). We currently support the following platforms:

|  ARCH/OS  | windows | linux | darwin |
| --------: | :-----: | :---: | :----: |
| **amd64** |   ✅    |  ✅  |   ✅   |
| **arm64** |   ✅    |  ✅  |   ✅   |
| **armv7** |   ❌    |  ✅  |   ❌   |

This means that LabStore runs on devices like the Raspberry Pi 2 and NAS units based on the armv7 architecture, as well as Apple Silicon and Intel-based Macs. If you architecture isn't currently supported, open an issue and we'll try to add it to our release workflow.


## Configuration

You configuration file for LabStore can be stored on the following locations, depending on the OS, and besides a local `labstore.yml` or `/etc/labstore/labstore.yml` when available:

| OS | Config Location |
| -- | --------------- |
| Linux | `~/.config/labstore/labstore.yml` |
| Mac | `~/Library/Application Support/labstore/labstore.yml` |
| Windows | `%APPDATA%/labstore/labstore.yml` |

We provide a configuration file example under [labstore.example.yml](labstore.example.yml). Any configuration setting can also be set using an environment variable, as well as a CLI argument, both following the same path structure. For example, `server.storage.data_dir` in the config file can be overridden by the `LABSTORE_SERVER_STORAGE_DATA_DIR` environment variable, or by the `--server-storage-data-dir` CLI argument.

In the future, we will provide a proper documentation web site. However, at this time, the easiest way to obtain any documentation about the configuration settings is to use the CLI help directly, namely for the `serve` command which uses nearly all configs:

```bash
labstore help serve
```

```
Run server for S3, IAM, and admin services

Usage:
  labstore serve [flags]

Flags:
      --admin-address-host string         Listening host for admin server (default "0.0.0.0")
      --admin-address-port string         Listening port for admin server (default "0.0.0.0")
      --admin-auth-access-key string      Administrator account access key (default "admin")
      --admin-auth-secret-key string      Administrator account secret key (default "admin")
  -h, --help                              help for serve
      --iam-address-host string           Listening host for IAM server (default "0.0.0.0")
      --iam-address-port string           Listening port for IAM server (default "0.0.0.0")
      --iam-db-max-idle-conns int         Maximum idle reader connections for the IAM database (default 3)
      --iam-db-max-open-conns int         Maximum open reader connections for the IAM database (default 3)
      --iam-db-read-cache-size-kib int    Cache size of each individual reader connection for the IAM database (default 65536)
      --iam-db-timeout-ms int             Connection timeout for the IAM database (default 5000)
      --iam-db-write-cache-size-kib int   Cache size of the writer connection for the IAM database (default 16384)
      --iam-db-write-chan-cap int         Buffered channel capacity for writing requests to the IAM database (default 32)
      --s3-address-host string            Listening host for S3-compatible server (default "0.0.0.0")
      --s3-address-port uint16            Listening port for S3-compatible server (default 6789)
      --s3-io-buffer-size int             Input/output buffer size in bytes (default 262144)
      --s3-paging-max-keys int            Hard limit for the maximum number of keys to return in paged requests (default 1000)
      --storage-data-dir string           Storage root path for objects and metadata (default "./data")
      --storage-keys-dir string           Storage root path for encryption keys (default "./data")
      --web-address-host string           Listening host for web ui server (default "0.0.0.0")
      --web-address-port uint16           Listening port for web ui server (default 6790)

Global Flags:
      --debug               Set debug level for logging
      --pprof               Enable profiler
      --pprof-host string   Profiler host (default "localhost")
      --pprof-port int      Profiler port (default 6060)
```


### Credentials

Client-side credentials are stored in the same configuration directory, alongside `labstore.yml` in a `credentials.yml` file, with the following format:

```yml
default_profile: local
profiles:
  local:
    access_key: admin
    secret_key: adminadmin
```

Each profile is identified by a key (e.g., `local`) under which we set an `access_key` and a `secret_key`. The top-level `default_profile` entry is set to a profile key (e.g., `local`), which means that, if no `--profile` CLI argument is used, the credentials in the default profile will be used.


## Features

- S3-compatible API.
  - Minimal support for basic API—single node, single drive only.
  - See [docs/api.md](./docs/api.md) for details.
  - Everything marked with 🟡 is supported, but we haven't ensured it fully meets the spec yet—e.g., it might be missing some response status codes, or not cover a particular use case or behavior.
  - As of now, there is no support for multipart upload, or any extra features around buckets or objects, like ACL, versioning, encryption, bucket-level policies, object locking, CORS configurations, etc.
- IAM-compatible API.
  - Minimal support for user and access key management (limited to one access key per user for now).
  - Minimal support for group management.
  - Minimal support for managed policies.
- CLI (one binary to rule them all).
- TUI (for S3 and IAM).
- Web UI (planned for `v0.2.0`).
  - Embedded in a single binary (default).
  - Downloaded in runtime (`go install` or custom builds only).


## Background

LabStore was born out of a perfect storm:

- Another project's open source solution entering [maintenance mode](https://github.com/minio/minio/issues/21647#issuecomment-3439134621), focusing on bug fixes and security patches exclusively.
- The need for content for [@DataLabTechTV](https://github.com/DataLabTechTV/)'s [YouTube channel](https://www.youtube.com/@DataLabTechTV).
- Curiosity to learn a new programming language (Go) and new tech stacks (SvelteKit).
- The desire to create a space for friends to learn, build, and push their portfolios to the next level.
- The age-old troublemaking question: how hard can it be?


### What LabStore Aims to Be

- An object store you can spin up quickly, locally, in a NAS, or in a home lab, to help you prototype a system.
- An object store with good integration with data tools.
- A "waste-no-time" local solution, with sane defaults and easy deployment.


### What LabStore Isn't Trying to Be

- A complete replacement for cloud solutions (use LabStore to prototype locally, but scale to AWS S3 or another compatible cloud service when your project grows).
- A project that starts as open source but then moves development behind closed doors, leaving users and contributors hanging. LabStore aims to stay transparent and community-friendly, so you can rely on it as it grows—the code lives here, no matter what else we do around the product.
- A corporate project. We’re volunteers—builders who create because we love it, not executives chasing the latest AI trend from an ivory tower.
