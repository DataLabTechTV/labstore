# Contributing

## Branch Naming

Please use the following naming conventions for branches. Remember that there are feature branches, but no fix branches or similar—those should be commits within a branch. Hotfix branches are to be created post-release, in emergencies, and release branches handle integration, testing and fixes pre-release, usually for a single project within the monorepo (e.g., `web/` or `backend/`). The stable version is always under the `main` branch.

### Feature branches

```
feature/<issue-number>-<short-description>
```

Example: `feature/123-add-login-page`

### Hotfix branches

```
hotfix/<short-description>
```

Example: `hotfix/fix-login-bug`

### Release branches

If multiple projects are combined in a single release, use the global unified version:

```
release/<version>
```

Example: `release/v1.0.0`

If a release is being done for a single project, then use the project-specific name and version:

```
release/<project>/<version>
```

Example: `release/web/v0.0.1`

## Commit Messages

We follow the conventional commits spec: https://www.conventionalcommits.org/en/v1.0.0/. And we specify a few of the optionals as well.

### Title and Body Format

Prefer all lowercase for title, unless absolutely required (e.g., uppercase env vars).

The body should be properly formatted text, but do not use markdown—while it works in GitHub, it can clutter messages from `git log`.

For a breaking change, always use the type with an `!`, as well as a body message with the `BREAKING CHANGE:` annotation. For example:

```
chore!: drop support for node 6

BREAKING CHANGE: use javascript features not available in node 6
```

### Type Scope

For the type scope, use nothing for top-level files (e.g., `justfile`, `.gitignore`, etc.), and use the project name (i.e., the name of the folder at the top-level, e.g., `web`, `backend`, etc.) as the optional scope.

An example for root-level:

```
chore: add node_modules to gitignore
```

An example for the `web/` project (web UI frontend):

```
feat(web): initialize svelte project
```

Or for the `backend/` project:

```
chore(backend): add logger dependency
```

## Repo Structure

We'll use a monorepo structure, where top-level directories correspond to independent projects.

```
monorepo/
├── .github/
├── web/
├── backend/
├── shared/
├── infra/
├── cli/
├── docs/
├── .gitignore
├── justfile
└── README.md
```

### Component Details

- `.github/` – GitHub templates (PRs, issues, etc.) and workflows (CI/CD)
- `web/` – web ui frontend
- `backend/` – REST API endpoint and admin tools
- `shared/` – shared assets, specs, etc. (no libs here—create another top-level project as those come along)
- `infra/` – CI/CD (to call from GH workflows), deployment scripts, Docker, etc.
- `cli/` – master command line tool that brings all projects together
- `docs/` – markdown documentation
- `.gitignore` – always use a single ignore file at the root, unless absolutely necessary

## Project Management

We'll use GitHub Projects for project management. This let's us organize issues and PRs in a Kanban board.

### Communication

- Use GitHub Discussions for new ideas and design questions.
- Use Issues to specify tasks (enhancements, bugs, etc.)
- Feature branches tie to issues through their ID
- Open PRs as soon as there is code to review (draft it early, make it visible).
- PRs must reference an issue.

### Workflow summary

1. Propose → in Discussions.
2. Track → Issue + Project board.
3. Implement → feature branch, PR, code review.
4. Merge → CI checks must pass (aspirational step).

Make it brief. Don't waste too much time writing up issues and PRs, but ensure all required information is there. Too much structure will slow us down. Too little structure will produce chaos. Be pragmatic.
