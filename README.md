# Microservice

A Go microservice starter built as a [Go workspace](go.work) with an HTTP gateway that talks to backend services over gRPC.

## Architecture

```
                 HTTP (JSON)
Client ─────────────────────►  Gateway  ─────── gRPC ──────►  Auth Service  ──────►  PostgreSQL
                               :6000                          :6001
```

- **Gateway** (`gateway/`) — public HTTP API. Validates and parses requests, forwards calls to the appropriate service over gRPC.
- **Auth Service** (`auth_service/`) — gRPC service handling user registration, credentials, and JWT issuing. Runs SQL migrations on startup.
- Shared libraries in `pkg/` (router, errors, database, config loader), shared protobuf messages in `general/`, shared config structs in `config/`.

## Tech stack

| Area       | Choice                                        |
|------------|-----------------------------------------------|
| Language   | Go ≥ 1.25                                     |
| RPC        | gRPC + Protocol Buffers                       |
| HTTP       | `net/http` with a custom router (`pkg/router`)|
| Database   | PostgreSQL (via `lib/pq` + `sql-migrate`)     |
| Auth       | bcrypt password hashing, HS256 JWT (`golang-jwt/jwt/v5`) |
| Deploy     | Docker / docker compose                       |

## Repository layout

```
├── gateway/          # HTTP gateway (public API)
├── auth_service/     # Auth microservice (gRPC) + migrations
├── pkg/              # Shared libraries: router, errors, database, yaml loader
├── config/           # Shared configuration structs
├── general/          # Shared protobuf messages
├── env.yaml.example  # Config template — copy to env.yaml
└── docker-compose.yml
```

## Getting started

### Prerequisites

- Go ≥ 1.25
- Docker & docker compose *(optional)*
- PostgreSQL *(only if running without docker)*

### 1. Configuration

```sh
cp env.yaml.example env.yaml   # then fill in real values
```

`env.yaml` is gitignored — never commit it.

### 2a. Run with docker compose

```sh
POSTGRES_PASSWORD=secret docker compose up --build -d
```

> For the full stack in docker, set `microservices.auth.ip` to `0.0.0.0`
> in `env.yaml`, and point the gateway at the service by hostname if you
> use container DNS names.

### 2b. Run locally

Start PostgreSQL, create a database named `authDB` (matching `env.yaml`),
then:

```sh
make run-auth      # terminal 1 — gRPC :6001
make run-gateway   # terminal 2 — HTTP :6000
```

Both services read `../env.yaml` relative to their directory.

## API endpoints

| Method | Path                   | Body                                                                                              | Success |
|--------|------------------------|---------------------------------------------------------------------------------------------------|---------|
| POST   | `/api/auth/register/`  | `{"username": "...", "password": "...", "password_confirm": "...", "email": "..."}`               | 201     |

Errors are returned as JSON:

```json
{"code": 400, "message": "username already taken"}
```

## Development

```sh
make build        # compile all services
make test         # run unit tests
make lint         # go vet across all modules
make tidy         # tidy all go.mod files
make proto        # regenerate protobuf/gRPC code
make help         # list all targets
```

Tests in `auth_service/service` use [`go-sqlmock`](https://github.com/DATA-DOG/go-sqlmock),
so they run without a database. See `CONTRIBUTING.md` for conventions.

## Security notes

- Secrets live only in `env.yaml` (gitignored). Use `env.yaml.example` as template.
- JWTs are signed HS256 with `secret_key` and expire after 24h.
- Passwords are hashed with bcrypt; they are never stored or logged in plaintext.
