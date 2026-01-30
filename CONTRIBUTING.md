# Contributing


## Etiquette

1. Introduce yourself in [Discussions](https://github.com/IllumiKnowLabs/labstore/discussions/categories/development) before opening your first PR.
    - Tell us a bit about yourself, and why you want to contribute.
    - We want to know who we're working on. After all this ir our baby, and it's still early.

2. Don't pickup issues that are already assigned to the team (check `Assignees`)
    - If we assigned it to us, then we're already working on it, or, at the very least, it's on our To Do list—we're using this opportunity to learn, so we don't want to delegate most of the feature issues.
    - If you still want to help with an already assigned issue, then post on [Discussions](https://github.com/IllumiKnowLabs/labstore/discussions/categories/development), tagging the issue and the assignee, so that we can coordinate and determine if you can help.
    - The most relevant issues you can pickup as an external contributor (i.e., from outside of the org) are the ones that address community bugs that do not require design changes.


## Repo Structure

We'll use a monorepo structure, where top-level directories correspond to independent projects.

```
labstore/
├── .github/
├── web/
├── client/
├── server/
├── cli/
├── cmd/
├── shared/
├── infra/
├── cli/
├── docs/
├── .gitignore
├── .dockerignore
├── .pre-commit-config.yml
├── justfile
├── Makefile
├── labstore.example.yml
└── README.md
```


### Component Details

- `.github/` – GitHub templates (PRs, issues, etc.), and actions and workflows (CI/CD).
- `web/` – Web UI SvelteKit project (static site).
- `client/` - S3 and IAM clients (library).
- `server/` – S3 and IAM API endpoints, and admin service (library).
- `cli/` – Central CLI tool for LabStore (server, client, TUI, etc.).
- `cmd/` – LabStore Cobra command (`labstore` binary).
- `infra/` – `Dockerfile`, `compose.yml`, scripts, etc.
- `shared/` – Shared assets, specs, etc. (no libraries here).
- `docs/` – Markdown documentation.
- `.gitignore` – Prefer a single ignore file at the root of the repo.
- `.dockerignore` – Files to be ignored during `COPY` on `Dockerfile` build.
- `.pre-commit-config.yml` – To setup linting—called by `just lint`, as a pre-commit hook, and during CI/CD.
- `justfile` – Project-wide tasks, largely related with manual testing, benchmarking, and infrastructure.
- `Makefile` – Helps with building backend and frontend, with correct `ldflags`, with or without embedded web assets, for multiple platforms, etc.
- `labstore.example.yml` – Example configuration file, with sane, usable defaults.
- `README.md` – Covers installation, configuration, features, and background story.


## Dev Requirements

### Install

If you're going to contribute to LabStore, you should install the following software (tested version and recommended installation instructions to follow):

[go](https://go.dev/doc/install) 1.25.5

```bash
# Download archive from link above, and then...
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin  # Persist this path in your shell.
go version
```

[node](https://nodejs.org/en/download) 25 w/ `npm`

```bash
curl -o- https://fnm.vercel.app/install | bash
fnm install 25
fnm use 25
node --version
npm --version
```

[just](https://github.com/casey/just?tab=readme-ov-file#installation) 1.43.0

```bash
npm install -g rust-just
just -V
```

[golangci-lint](https://golangci-lint.run/docs/welcome/install/local/) 2.8.0

```bash
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.8.0
golangci-lint --version
```

[pre-commit](https://pre-commit.com/#install) 4.2.0

```bash
sudo apt install -y pre-commit
pre-commit --version
```

[jq](https://jqlang.org/download/) 1.7.1

```bash
sudo apt install -y jq
jq --version
```

[yq](https://github.com/mikefarah/yq/#install) 4.45.4

```bash
go install github.com/mikefarah/yq/v4@latest
yq --version
```

[mc](https://github.com/minio/mc?tab=readme-ov-file#install-from-source) RELEASE.2025-08-13T08-35-41Z

```bash
go install github.com/minio/mc@latest
mc --version  # DEVELOPMENT.GOGET regardless
```

[rclone](https://rclone.org/downloads/) 1.60.1

```bash
sudo apt install -y rclone
rclone version
```

[s5cmd](https://github.com/peak/s5cmd?tab=readme-ov-file#installation) 2.3.0

```bash
go install github.com/peak/s5cmd/v2@master
s5cmd version  # v0.0.0-dev regardless
```

[docker](https://docs.docker.com/engine/install/debian/) 29.1.5

```bash
# Debian

# Uninstall old versions:
sudo apt remove $(dpkg --get-selections docker.io docker-compose docker-doc podman-docker containerd runc | cut -f1)

# Add Docker's official GPG key:
sudo apt update
sudo apt install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/debian
Suites: $(. /etc/os-release && echo "$VERSION_CODENAME")
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF

sudo apt update

sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

docker version
```

[iperf3](https://iperf.fr/iperf-download.php) 3.18

```bash
sudo apt install -y iperf3
```

[warp](https://github.com/minio/warp) 1.3.1

```bash
go install github.com/minio/warp@latest
warp --version  # (dev) - (dev) regardless
```

[duckdb](https://duckdb.org/install/?platform=linux&environment=cli) 1.4.3

```bash
curl https://install.duckdb.org | sh
duckdb --version
```

[termgraph](https://github.com/mkaz/termgraph) 0.7.4

```bash
uv tool install termgraph==0.7.4
```

[sqlite3](https://www.sqlite.org/download.html) 3.46.1

```bash
sudo apt install -y sqlite3
sqlite3 --version
```

[http](https://httpie.io/docs/cli/installation)

```bash
curl -SsL https://packages.httpie.io/deb/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/httpie.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/httpie.gpg] https://packages.httpie.io/deb ./" | sudo tee /etc/apt/sources.list.d/httpie.list > /dev/null
sudo apt update && sudo apt install -y
```


### Setup

```bash
just install-hooks
```

You can run `just lint` to make sure `pre-commit` and `golangci-lint` are working as expected.


#### vscode

If you're using vscode, here's a quick `launch.json` example to help with debugging:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "server",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "cwd": "${workspaceFolder}",
      "program": "${workspaceFolder}/cmd/labstore/",
      "env": {
        // "LABSTORE_BACKEND_S3_SERVER_HOST": "localhost",
      },
      "args": [
        "serve",
        "--debug",
        // "--admin-auth-access-key", "root",
      ],
      "console": "integratedTerminal",
    },
}
```


### Build

The provided `Makefile` contains all build targets we require. See the table below for a summary:

| Target | Description |
| ------ | ----------- |
| `all` | Same as `build`. |
| `assets` | Build `web` and copy `web/build/` to `server/router/assets/`. |
| `cli` | Build `assets` and `labstore` binary. |
| `web` | Build web UI static site (assets for Go `server`). |
| `build-<OS>-<arch>` | Cross-compilation for the corresponding OS and arch (binary under `bin/labstore-<OS>)-<arch>`, ending in `.exe` for Windows. |
| `build` | Compile for the current platform. |
| `build-all` | Cross-compile for all supported platforms. |
| `build-runtime` | Compile for the current platform, without embedded web UI assets. |
| `run` | Run server and web UI, for testing during development. |
| `profile` | Run the Go profiler (`pprof`), capturing for 60 seconds, and opening the web UI for analysis. |
| `test` | Run all tests (currently only Go is setup). |
| `clean-assets` | Delete the `assets/` directory under `server/router/`. |
| `clean-cli` | Run `clean-assets` and delete `bin/`. |
| `clean-web` | Delete node build under `web/`. |
| `clean` | Run `clean-cli` and `clean-web`. |
| `dist-clean` | Run `clean` and delete debug binaries and Go cache. |


## Git

### Branch Naming

The stable version is always under the `main` branch. Please use the following naming convention for any other branches, which is the default behavior for GitHub Projects.

```
<issue-number>-<issue-title>
```

Example: `123-add-login-page`


### Tagging

We support `testing-*` tags for internal use, as well as version tags in the format `v*.*.*` (e.g., `v0.1.0`, `v0.1.0-alpha.1`, `v0.2.0-beta.10`, `v0.2.0-rc.1`).

When creating a release tag, make sure you create annotated tags, if you want `ldflags` to be set correctly. For example:

```bash
git tag -a v0.1.0-alpha.10 -m "release: v0.1.0-alpha.10
```

Any non-annotated tag will still trigger the `docker` and `release` workflows, but the produced binaries, both inside the Docker image and the GitHub Release won't display the correct information when running `labstore version`.


### Commit Messages

We follow the conventional commits spec: https://www.conventionalcommits.org/en/v1.0.0/. And we specify a few of the optionals as well.


#### Title and Body Format

Prefer all lowercase for title, unless absolutely required (e.g., uppercase env vars).

The body should be properly formatted text, but do not use markdown—while it works in GitHub, it can clutter messages from `git log`.

For a breaking change, always use the type with an `!`, as well as a body message with the `BREAKING CHANGE:` annotation. For example:

```
chore!: drop support for node 6

BREAKING CHANGE: use javascript features not available in node 6
```


#### Type Scope

For the type scope, use nothing for top-level files (e.g., `justfile`, `.gitignore`, etc.), but only when there is no other option. Otherwise, always use the project name (i.e., the name of the folder at the top-level, e.g., `web`, `server`, etc.) as the optional scope.

An example for root-level:

```
chore: add node_modules to gitignore
```

Another example for root-level, affecting the `justfile` for the `server` project:

```
chore(server): add server run command
```

An example for the `web/` project (web UI frontend):

```
feat(web): initialize svelte project
```

Or for the `server/` project:

```
chore(server): add logger dependency
```


## CI/CD

We define one GitHub Action and four workflows:

```
.github/
├── actions
│   └── setup
│       └── action.yml
└── workflows
    ├── lint.yml
    ├── test.yml
    ├── docker.yml
    └── release.yml
```

| GHA | Description |
| -------- | ----------- |
| `setup` | Used by `lint`, `test`, and `release`, to setup Go, Node, building `web` and copying the output to the assets directory, required for Go building, linting, and testing. |
| `lint` | Run linting using `pre-commit`, on PRs for `main`, `backend`, `web`, and `release/**`. |
| `test` | Run all tests (currently only Go is setup), on PRs for `main`, `backend`, `web`, and `release/**`. |
| `docker` | Build and publish the Docker image to GHCR, on a push to `main` or tags `v*.*.*` and `testing-*`. |
| `release` | Builds the binaries and the web assets, generating a changelog and a GitHub release with all of those, on a push to tags `v*.*.*` and `testing-*`.

## Project Management

We use GitHub Projects for project management. This lets us organize issues and PRs in a Kanban board.


### Communication

- Use GitHub Discussions for new ideas and design questions.
- Use Issues to specify tasks (`enhancement`, `bug`, etc.) for a given project (e.g., `cli`, `server`, `web`, etc.).
- Feature branches tie to issues through their ID
- Open PRs as soon as there is code to review (draft it early, make it visible).
- PRs must reference an issue.


### Workflow summary

1. Propose → in discussions.
2. Track → issue + project board.
3. Implement → feature branch, PR, code reviewing.
4. Merge → CI checks must pass (mandatory for `main`).

Make it brief. Don't waste too much time writing up issues and PRs, but ensure all required information is there. Too much structure will slow us down. Too little structure will produce chaos. Be pragmatic.
