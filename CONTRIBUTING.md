# Contributing

Thanks for your interest in contributing!

## Setup

1. Install Go >= 1.25.
2. `cp env.yaml.example env.yaml` and fill in local values (never commit this file).
3. `make test` should pass before and after your changes.

## Workflow

- Branch off `develop`; open PRs against `develop`.
- Keep commits small and descriptive: `area: what changed` (e.g. `gateway: add login endpoint`).
- Run `make lint && make test` before pushing.

## Conventions

- Each top-level folder is its own Go module; the repo root is a Go workspace.
- New HTTP routes go in `<service>/routes/http.go`; handlers in `<service>/handlers/`.
- Return client-facing errors via `pkg/errors.New(code, errMsg, responseMessage)` — the panic-recovery middleware converts them to JSON responses.
- Database changes require a new migration file in `auth_service/migrations/<env>/NNN.description.sql` (PostgreSQL syntax).
- Regenerate protobuf code with `make proto` after editing any `.proto` file; commit both `.proto` and generated files.
